package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"order/internal/storage"
	"os"
	"time"
)

type Producer struct {
	P *kafka.Producer
}

func (pr *Producer) Run(ctx context.Context, topic string, limit int) {
	connStr, exists := os.LookupEnv("DB_CONNECTION")
	if !exists {
		panic("No connection sting to db")
	}
	s, err := storage.InitStorage(connStr)
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
	fmt.Println(answer.String())

}

func InitProducer() (*Producer, error) {
	bootstrapServer, bsExists := os.LookupEnv("KAFKA_BOOTSTRAP_SERVER")
	group, groupExists := os.LookupEnv("KAFKA_GROUP_ID")
	reset, resetExists := os.LookupEnv("KAFKA_RESET")

	if !bsExists || !groupExists || !resetExists {
		panic("No kafka config")
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          group,
		"auto.offset.reset": reset,
	})

	if err != nil {
		return nil, err
	}

	return &Producer{
		P: producer,
	}, nil
}
