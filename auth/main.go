package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/pnaskardev/URL-Shortner-V1/auth/config"
	"github.com/pnaskardev/URL-Shortner-V1/auth/infrastructure/container"
)

// THIS WILL BE PURELY AN RPC SERVER AND NOTHING ELSE
func main() {

	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// config := config.GetConfig()

	rpcCont := container.Init()

	go func() {
		address := rpcCont.Server.GetAddr()
		lis, err := net.Listen("tcp", address)
		if err != nil {
			panic(err)
		}

		log.Printf("RPC server listening on %s", address)

		if err := rpcCont.Server.GetServer().Serve(lis); err != nil {
			panic(err)
		}
	}()
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
