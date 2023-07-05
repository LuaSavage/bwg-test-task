package msgbroker

import (
	"github.com/LuaSavage/bwg-test-task/service-b/pkg/client/msgbroker/dto"
	"github.com/LuaSavage/bwg-test-task/service-b/pkg/logging"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	Consumer *kafka.Consumer
	Logger   logging.Logger
}

func NewKafkaConsumer(dto *dto.NewConsumerDTO) (*KafkaConsumer, error) {
	// Kafka consumer configuration
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  dto.BrokerAdress, // Kafka broker address
		"group.id":           dto.GroupId,
		"auto.offset.reset":  dto.AutoOffsetReset,
		"enable.auto.commit": dto.EnableAutoCommit,
	})
	if err != nil {
		return nil, err
	}

	// Subscribe to the reply topic
	err = consumer.SubscribeTopics([]string{dto.Topic}, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		Consumer: consumer,
		Logger:   logging.GetLogger(),
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