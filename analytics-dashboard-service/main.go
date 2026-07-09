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
	"github.com/pnaskardev/URL-Shortner-V1/analytics-dashboard-service/config"
	"github.com/pnaskardev/URL-Shortner-V1/analytics-dashboard-service/infrastructure/database"
	"github.com/pnaskardev/URL-Shortner-V1/analytics-dashboard-service/infrastructure/queue"
)

func main() {

	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	config := config.GetConfig()

	_ = database.ConnectToPostgres()

	q, err := queue.NewQueueClient()
	if err != nil {
		panic(err)
	}

	err = q.DeclareQueue()
	if err != nil {
		slog.Error("RABBIT MQ ERROR", "ERROR", err)
		panic(err)
	}

	app := fiber.New(
		fiber.Config{
			AppName: "ANALYTICS/DASHBOARD SERVICE",
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
	// q.Close()
	fmt.Println("Fiber was successful shutdown.")

}
