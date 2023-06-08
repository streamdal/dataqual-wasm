package dataqual_wasm

import (
	"strings"

	"github.com/pkg/errors"
	sarama "github.com/streamdal/shopify-sarama"
)

type SaramaKafka struct {
	Client   sarama.Client
	Producer sarama.SyncProducer
	Topic    string
}

func setupSarama(topic string) (*SaramaKafka, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_6_0_0

	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true

	client, err := sarama.NewClient([]string{"localhost:9092"}, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create sarama client")
	}

	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create sarama cluster admin")
	}

	if err := admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return nil, errors.Wrap(err, "unable to create kafka topic")
		}
	}

	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}

	return &SaramaKafka{
		Client:   client,
		Producer: p,
		Topic:    topic,
	}, nil
}

func (k *SaramaKafka) Cleanup() {
	// TODO: delete topic
	k.Producer.Close()
	k.Client.Close()
}
