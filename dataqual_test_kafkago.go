package dataqual_wasm

import (
	"context"
	"fmt"
	"strings"
	"time"

	kafka "github.com/streamdal/segmentio-kafka-go"
)

const (
	KafkaAddress = "localhost:9092"
	KafkaTopic   = "dqtest1"
)

type Segmentio struct {
	Conn   *kafka.Conn
	Topic  string
	Writer *kafka.Writer
}

func setupSegmentio(topic string) (*Segmentio, error) {
	dialer := &kafka.Dialer{
		Timeout: time.Second * 5,
	}

	conn, err := dialer.DialContext(context.Background(), "tcp", KafkaAddress)
	if err != nil {
		return nil, fmt.Errorf("unable to dial '%s': %s", KafkaAddress, err)
	}

	if err := conn.CreateTopics(kafka.TopicConfig{Topic: KafkaTopic, NumPartitions: 1, ReplicationFactor: 1}); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			panic(err)
		}
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{KafkaAddress},
		BatchSize:    1,
		RequiredAcks: 1,
	})

	return &Segmentio{
		Conn:   conn,
		Topic:  topic,
		Writer: w,
	}, nil
}

func (s *Segmentio) Cleanup() {
	s.Conn.Close()
}
