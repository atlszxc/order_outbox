package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"order/internal/kafka"
	"order/internal/order/handler"
	"order/internal/server"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	port, portExist := os.LookupEnv("PORT")
	if !portExist {
		os.Exit(1)
	}

	s := server.InitServer(port, server.DEBUG)

	consumer, conErr := kafka.InitConsumer()
	if conErr != nil {
		fmt.Println(conErr)
		os.Exit(1)
	}

	producer, kafkaErr := kafka.InitProducer()
	if kafkaErr != nil {
		fmt.Println(kafkaErr)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	conCtx, conCancel := context.WithCancel(context.Background())

	go producer.Run(ctx, "order", 50)
	go consumer.Listen(conCtx)

	s.AddRoute(server.POST, "/create", handler.CreateOrderHandler)
	s.AddRoute(server.PATCH, "/update", handler.UpdateStatusHandler)

	err := s.Listen()
	if err != nil {
		fmt.Println(err)
		cancel()
		conCancel()
		os.Exit(1)
	}
}
