package queue

import (
	"context"
	"fmt"
	"sync"

	"github.com/pnaskardev/URL-Shortner-V1/analytics-dashboard-service/helpers/constants"
	amqp "github.com/rabbitmq/amqp091-go"
)

// pendingConfirm holds the result channel for a single in-flight publish.
type pendingConfirm struct {
	done chan bool // receives ack=true or nack=false; closed on channel failure
}

type QueueClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	mu       sync.Mutex
	pending  map[uint64]*pendingConfirm
	nextTag  uint64 // mirrors the broker's delivery tag counter
	closed   bool
	closeErr error

	confirms chan amqp.Confirmation
	closeCh  chan *amqp.Error
	doneCh   chan struct{} // closed when the confirm loop exits
}

func NewQueueClient() (*QueueClient, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	if err := ch.Confirm(false); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}

	c := &QueueClient{
		conn:     conn,
		channel:  ch,
		pending:  make(map[uint64]*pendingConfirm),
		confirms: ch.NotifyPublish(make(chan amqp.Confirmation, 256)),
		closeCh:  ch.NotifyClose(make(chan *amqp.Error, 1)),
		doneCh:   make(chan struct{}),
	}

	go c.confirmLoop()
	return c, nil
}

// confirmLoop drains confirmations and resolves pending publishes.
// Exits when the channel closes; at that point every still-pending publish
// is failed via a closed done channel.
func (c *QueueClient) confirmLoop() {
	defer close(c.doneCh)

	for {
		select {
		case conf, ok := <-c.confirms:
			if !ok {
				c.failAllPending()
				return
			}
			c.mu.Lock()
			p, found := c.pending[conf.DeliveryTag]
			if found {
				delete(c.pending, conf.DeliveryTag)
			}
			c.mu.Unlock()
			if found {
				p.done <- conf.Ack
			}

		case err := <-c.closeCh:
			c.mu.Lock()
			c.closed = true
			if err != nil {
				c.closeErr = err
			} else {
				c.closeErr = fmt.Errorf("channel closed")
			}
			c.mu.Unlock()
			c.failAllPending()
			return
		}
	}
}

func (c *QueueClient) failAllPending() {
	c.mu.Lock()
	for tag, p := range c.pending {
		close(p.done) // closed channel signals failure (no value sent)
		delete(c.pending, tag)
	}
	c.mu.Unlock()
}

func (c *QueueClient) DeclareQueue() error {
	if err := c.channel.ExchangeDeclare(constants.URL_EXCHANGE, "direct", true, false, false, false, nil); err != nil {
		return err
	}

	if _, err := c.channel.QueueDeclare(
		constants.URL_CREATED_QUEUE,
		true,
		false,
		false,
		false,
		amqp.Table{amqp.QueueTypeArg: amqp.QueueTypeQuorum},
	); err != nil {
		return err
	}

	return c.channel.QueueBind(
		constants.URL_CREATED_QUEUE,
		constants.URL_CREATED_QUEUE,
		constants.URL_EXCHANGE,
		false,
		nil,
	)
}

// Publish sends a message and blocks until the broker confirms it, the context
// expires, or the channel dies. Multiple goroutines can call Publish concurrently.
func (c *QueueClient) Publish(ctx context.Context, queueName string, body []byte) error {
	p := &pendingConfirm{done: make(chan bool, 1)}

	// Register BEFORE publishing so we can't miss a confirm that arrives
	// between the publish call and the map insert.
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return fmt.Errorf("queue client closed: %w", c.closeErr)
	}
	c.nextTag++
	tag := c.nextTag
	c.pending[tag] = p
	c.mu.Unlock()

	err := c.channel.PublishWithContext(
		ctx,
		constants.URL_EXCHANGE,
		queueName,
		false, false,
		amqp.Publishing{
			ContentType:  "application/x-protobuf",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		c.mu.Lock()
		delete(c.pending, tag)
		c.mu.Unlock()
		return err
	}

	select {
	case ack, ok := <-p.done:
		if !ok {
			return fmt.Errorf("channel closed before confirm")
		}
		if !ack {
			return fmt.Errorf("message nacked by broker")
		}
		return nil

	case <-ctx.Done():
		// Leave the entry in pending; the confirm loop will clean it up
		// when the confirm eventually arrives (or when the channel dies).
		// We just stop waiting.
		return ctx.Err()
	}
}

func (c *QueueClient) Consume(queuename string) (<-chan amqp.Delivery, error) {
	msgs, err := c.channel.Consume(
		queuename, "",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		return nil, err
	}

	return msgs, nil

}

func (c *QueueClient) Close() {
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
	<-c.doneCh // wait for confirm loop to drain
}
