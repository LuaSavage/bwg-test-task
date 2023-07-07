package msgbroker

import (
	"fmt"

	"github.com/LuaSavage/bwg-test-task/service-b/internal/config"
	"github.com/LuaSavage/bwg-test-task/service-b/pkg/logging"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	Consumer *kafka.Consumer
	Logger   logging.Logger
}

func NewKafkaConsumer(cfg config.KafkaConfig, logger logging.Logger) (*KafkaConsumer, error) {
	// Kafka consumer configuration
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), // Kafka broker address
		"group.id":           "service_b_consumer_group",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "false",
	})
	if err != nil {
		return nil, err
	}

	// Subscribe to the reply topic
	err = consumer.SubscribeTopics([]string{cfg.Topic}, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		Consumer: consumer,
		Logger:   logger,
	}, nil
}

func (k *KafkaConsumer) ConsumeMessage() (response []byte) {
	var message *kafka.Message
	defer k.Consumer.CommitMessage(message)
	for {
		// Poll for Kafka messages
		message, err := k.Consumer.ReadMessage(-1)
		if err != nil {
			k.Logger.Errorf("Failed to read message: %v\n", err)
			continue
		}

		response = message.Value
		return
	}
}

/*
func (k *KafkaConsumer) ConsumeMessageWithKey(messageKey string) (response []byte) {
	var message *kafka.Message
	defer k.Consumer.CommitMessage(message)
	for {
		// Poll for Kafka messages
		message, err := k.Consumer.ReadMessage(-1)
		if err != nil {
			k.Logger.Errorf("Failed to read message: %v\n", err)
			continue
		}

		// Process the reply
		if string(message.Key) == messageKey {
			response = message.Value
			return
		}
	}
}
*/

func (k *KafkaConsumer) Close() error {
	return k.Consumer.Close()
}
