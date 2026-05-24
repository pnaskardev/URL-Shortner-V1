package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/google/uuid"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/api/routes"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/config"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/infrastructure/database"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/infrastructure/queue"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	config := config.GetConfig()

	dbClient := database.ConnectToPostgres()

	q, err := queue.NewQueueClient()
	if err != nil {
		panic(err)
	}

	q.DeclareQueue("url_created_queue")

	customLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))
	slog.SetDefault(customLogger)

	app := fiber.New(
		fiber.Config{
			AppName: "SHORTENER SERVICE",
		},
	)
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(func(c fiber.Ctx) error {

		txID := c.Get("X-Transaction-ID")
		if txID == "" {
			txID = uuid.NewString()
		}
		c.Locals("transaction_id", txID)
		c.Set("X-Transaction-ID", txID)

		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		c.Append("Server-Timing", "app;dur="+duration.String())
		return err
	})

	routes.ApiRouter(app, dbClient, q)
	// requestClient := &httpClients.RetryableHTTPClient{
	// 	Client: &http.Client{
	// 		Timeout: 5 * time.Second,
	// 	},
	// 	Config: httpClients.RetryConfig{
	// 		MaxRetries: 3,
	// 		BaseDelay:  100 * time.Millisecond,
	// 		MaxDelay:   2 * time.Second,
	// 	},
	// }

	// Register all of the routes
	// routes.ApiRouter(app, requestClient, dbClient)

	port := ":" + config.Port

	go func() {
		if err := app.Listen(port, fiber.ListenConfig{
			EnablePrefork:     false,
			EnablePrintRoutes: true,
		}); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")

	// Your cleanup tasks go here
	// db.Close()
	// redisConn.Close()
	fmt.Println("Fiber was successful shutdown.")

}
