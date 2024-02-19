package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"order/internal/kafka/dtos"
	"order/internal/storage"
	"os"
	"time"
)

type Consumer struct {
	C *kafka.Consumer
}

func (c *Consumer) Listen(ctx context.Context) {
	connStr, exists := os.LookupEnv("DB_CONNECTION")
	if !exists {
		panic("No connection sting to db")
	}
	s, storageErr := storage.InitStorage(connStr)
	if storageErr != nil {
		panic(storageErr)
	}

	err := c.C.SubscribeTopics([]string{"confirmOrder"}, nil)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.C.ReadMessage(5 * time.Second)
			if err == nil {
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
				var res dtos.ResponseDeliveredStatusDto
				convertErr := json.Unmarshal(msg.Value, &res)
				if convertErr != nil {
					fmt.Println(convertErr)
				}
				s.UpdateDeliveredStatus(res.OrderId, res.Status)
			} else if !err.(kafka.Error).IsTimeout() {
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}

	c.C.Close()
}

func InitConsumer() (*Consumer, error) {
	bootstrapServer, bsExists := os.LookupEnv("KAFKA_BOOTSTRAP_SERVER")
	group, groupExists := os.LookupEnv("KAFKA_GROUP_ID")
	reset, resetExists := os.LookupEnv("KAFKA_RESET")

	if !bsExists || !groupExists || !resetExists {
		panic("No kafka config")
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          group,
		"auto.offset.reset": reset,
	})

	if err != nil {
		return nil, err
	}

	return &Consumer{
		C: c,
	}, nil
}
