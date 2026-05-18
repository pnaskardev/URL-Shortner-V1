package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/pnaskardev/URL-Shortner-V1/core/api/routes"
	"github.com/pnaskardev/URL-Shortner-V1/core/config"
	requesthelper "github.com/pnaskardev/URL-Shortner-V1/core/helpers/requestHelper"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	config := config.GetConfig()

	customLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))
	slog.SetDefault(customLogger)

	app := fiber.New(
		fiber.Config{
			AppName: "CORE",
		},
	)
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		c.Append("Server-Timing", "app;dur="+duration.String())
		return err
	})

	requestClient := &requesthelper.RetryableHTTPClient{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		Config: requesthelper.RetryConfig{
			MaxRetries: 3,
			BaseDelay:  100 * time.Millisecond,
			MaxDelay:   2 * time.Second,
		},
	}

	// Register all of the routes
	routes.ApiRouter(app, requestClient)

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
