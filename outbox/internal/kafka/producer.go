package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"os"
	"outbox/internal/storage"
	"time"
)

type Producer struct {
	P *kafka.Producer
}

func (pr *Producer) Run(ctx context.Context, topic string) {
	connStr, exists := os.LookupEnv("DB_CONNECTION")
	if !exists {
		panic("No connection sting to db")
	}
	s := storage.GetStorage(connStr)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			ids, err := s.GetOutbox(50)
			if err != nil {
				panic(err)
			}

			for _, id := range ids {
				pr.SendMessage(topic, id)
				s.UpdateStatus(id)
			}

			time.Sleep(5 * time.Second)
		}
	}
}

func (pr *Producer) SendMessage(topic string, msg any) {
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
