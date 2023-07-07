package msgbroker

import (
	"fmt"

	"github.com/LuaSavage/bwg-test-task/service-a/internal/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/echo/v4"
)

type KafkaProducer struct {
	Producer *kafka.Producer
	Topic    *string
	Logger   echo.Logger
}

func NewKafkaProducer(cfg config.KafkaConfig, logger echo.Logger) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	})
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		Producer: producer,
		Topic:    &cfg.Topic,
		Logger:   logger,
	}, nil
}

func (k *KafkaProducer) ProduceMessage(messageKey []byte, messageValue []byte) error {
	// Produce the Kafka message
	err := k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: k.Topic, Partition: kafka.PartitionAny},
		Key:            messageKey,
		Value:          messageValue,
	}, nil)
	if err != nil {
		return err
	}

	// Wait for message delivery
	deliveryReport := <-k.Producer.Events()
	message := deliveryReport.(*kafka.Message)

	if message.TopicPartition.Error != nil {
		err = fmt.Errorf("Failed to deliver message: %v\n", message.TopicPartition.Error)
		k.Logger.Error(err.Error())
		return err
	}

	k.Logger.Infof("Message delivered to topic: %s, partition: %d, offset: %d\n",
		*message.TopicPartition.Topic, message.TopicPartition.Partition, message.TopicPartition.Offset)
	return nil
}

func (k *KafkaProducer) Close() {
	k.Producer.Close()
}
