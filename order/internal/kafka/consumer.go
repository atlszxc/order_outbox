package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"order/internal/storage"
	"time"
)

type ResponseDeliveredStatusDto struct {
	OrderId int
	Status  bool
}

type Consumer struct {
	C *kafka.Consumer
}

func (c *Consumer) Listen(ctx context.Context) {
	s, storageErr := storage.InitStorage("postgres://root:postgres@orderdb:5432/order")
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
				var res ResponseDeliveredStatusDto
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
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "order",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	return &Consumer{
		C: c,
	}, nil
}
