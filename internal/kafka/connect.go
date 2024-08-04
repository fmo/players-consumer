package kafka

import (
	"crypto/tls"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"os"
)

type Kafka struct {
	Reader *kafka.Reader
}

func NewKafka() Kafka {
	return Kafka{
		Reader: connect(os.Getenv("KAFKA_TOPIC"), os.Getenv("KAFKA_GROUP")),
	}
}

func connect(topic string, groupId string) *kafka.Reader {
	mechanism, _ := scram.Mechanism(
		scram.SHA256,
		os.Getenv("KAFKA_USERNAME"),
		os.Getenv("KAFKA_PASSWORD"),
	)

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_HOST")},
		GroupID: groupId,
		Topic:   topic,
		Dialer: &kafka.Dialer{
			SASLMechanism: mechanism,
			TLS:           &tls.Config{},
		},
	})
}
