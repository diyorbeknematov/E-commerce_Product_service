package consumer

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer interface {
	ConsumeMessages(handler func(message []byte)) error
	Close()
}

type ConsumerKafkaImpl struct {
	reader *kafka.Reader
	logger *slog.Logger
}

func NewConsumerKafka(brokerAddrs []string, topic string, groupID string, logger *slog.Logger) KafkaConsumer {
	consumer := &ConsumerKafkaImpl{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokerAddrs,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
		logger: logger,
	}
	logger.Info("kafka consumer created", "topic", consumer.reader.Config().Topic)
	return consumer
}

func (c *ConsumerKafkaImpl) ConsumeMessages(handler func(message []byte)) error {
	c.logger.Info("kafka consumer started consuming messages", "topic", c.reader.Config().Topic)
	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}
		handler(m.Value)
	}
}

func (c *ConsumerKafkaImpl) Close() {
	c.reader.Close()
}
