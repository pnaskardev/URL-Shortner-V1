package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pnaskardev/URL-Shortner-V1/auth/config"
)

// THIS WILL BE PURELY AN RPC SERVER AND NOTHING ELSE
func main() {

	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// config := config.GetConfig()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	fmt.Println("Running cleanup tasks...")

	// Your cleanup tasks go here
	// db.Close()
	// redisConn.Close()
	fmt.Println("Fiber was successful shutdown.")

}
