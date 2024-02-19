package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"order/internal/storage"
	"time"
)

type Producer struct {
	P *kafka.Producer
}

func (pr *Producer) Run(ctx context.Context, topic string, limit int) {
	s, err := storage.InitStorage("postgres://root:postgres@orderdb:5432/order")
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			ids, storageErr := s.GetCompletedOrderId(limit)
			if storageErr != nil {
				panic(storageErr)
			}

			for _, id := range ids {
				pr.SendMessage(topic, id)
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func (pr *Producer) SendMessage(topic string, msg int) {
	v, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	delChan := make(chan kafka.Event)

	produceErr := pr.P.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          v,
	}, delChan)

	if produceErr != nil {
		fmt.Println(err)
	}

	answer := <-delChan
	z := answer.(*kafka.Message)
	fmt.Println("kafka")
	fmt.Println(z.Value)
}

func InitProducer() (*Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "order",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	return &Producer{
		P: producer,
	}, nil
}
