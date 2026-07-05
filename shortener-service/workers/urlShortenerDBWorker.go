package workers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	urlshortenerv1 "github.com/pnaskardev/URL-Shortner-V1/contracts/gen/go/urlshortener/v1"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/helpers/constants"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/infrastructure/database/models"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/infrastructure/queue"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	StartShortenerService()
}

type repository struct {
	dbClient    *gorm.DB
	queueClient *queue.QueueClient
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewShortenerIngestor(ctx context.Context, dbClient *gorm.DB, queueClient *queue.QueueClient) Repository {
	ctx, cancelFunc := context.WithCancel(ctx)
	return &repository{
		dbClient:    dbClient,
		queueClient: queueClient,
		ctx:         ctx,
		cancel:      cancelFunc,
	}
}

func (r *repository) StartShortenerService() {

	// Separate Go Routine for the ingestion service
	go func() {

		// Infinite Loop for the long process
		for {
			if err := r.runStartShortenerService(); err != nil {
				slog.Error("ERROR", "ShortenerService Crashed", err)
			}

			// If the process is done it will return else after 2 seconds it will restart
			select {
			case <-r.ctx.Done():
				return
			case <-time.After(2 * time.Second):
			}
		}
	}()

}

func (r *repository) runStartShortenerService() error {

	msgs, err := r.queueClient.Consume(constants.URL_CREATED_QUEUE)
	if err != nil {
		return err
	}

	for {
		select {
		case <-r.ctx.Done():
			return nil

		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("delivery channel closed")
			}

			event := new(urlshortenerv1.UrlCreatedEvent)
			if err := proto.Unmarshal(msg.Body, event); err != nil {
				slog.Error("bad event, dropping", "error", err)
				msg.Nack(false, false) // don't requeue a message that will never parse
				continue
			}

			fmt.Printf("%+v\n", event) // whole struct
			// // or, more readable:
			// // b, _ := protojson.Marshal(event); fmt.Println(string(b))

			shortenedURL := models.ShortenedURL{
				ShortURLKey: event.GetShortenedUrlKey(),
				LongURL:     event.GetLongUrl(),
				UserID:      event.GetUserId(),
			}

			err := r.dbClient.Transaction(func(tx *gorm.DB) error {
				event_id := event.GetId()
				res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.ProcessedEvents{
					EventID: event_id,
				})

				if res.Error != nil {
					return res.Error
				}

				// already processed
				// none of the rows were created that means this id already existied and we can return
				if res.RowsAffected == 0 {
					return nil
				}

				// Nothing existed and we need to put this row in
				return tx.Create(&shortenedURL).Error
			})

			if err != nil {
				slog.Error("persist failed, requeue", "error", err)
				msg.Nack(false, true) // requeue → retry
				continue
			}

			// false because I dont want to ACK all of the messages before this
			// I only want to ACK the current one
			msg.Ack(false)

		}

	}

}
