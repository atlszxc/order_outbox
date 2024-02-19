package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"outbox/internal/kafka"
	"outbox/internal/server"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	s := server.InitServer(":8081", server.DEBUG)
	c, consumerErr := kafka.InitConsumer()
	if consumerErr != nil {
		fmt.Println(consumerErr)
		os.Exit(1)
	}

	p, producerErr := kafka.InitProducer()
	if producerErr != nil {
		fmt.Println(producerErr)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	pCtx, pCancel := context.WithCancel(context.Background())

	go c.Listen(ctx)
	go p.Run(pCtx, "delivery")

	err := s.Listen()
	if err != nil {
		fmt.Println(err)
		cancel()
		pCancel()
		os.Exit(1)
	}
}
