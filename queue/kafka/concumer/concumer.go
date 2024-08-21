package consumer

import (
	"context"
	"log"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer interface {
	ConsumeMessages(ctx context.Context, handler func(message []byte)) error
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
	log.Println("kafka consumer created", "topic", consumer.reader.Config().Topic)
	logger.Info("kafka consumer created", "topic", consumer.reader.Config().Topic)
	return consumer
}

func (k *ConsumerKafkaImpl) ConsumeMessages(ctx context.Context, handler func(message []byte)) error {
	defer k.Close()

	for {
		select {
		case <-ctx.Done():
			k.logger.Info("Consumer shutting down")
			return nil
		default:
			m, err := k.reader.ReadMessage(ctx)
			if err != nil {
				k.logger.Error("Error reading message", "error", err)
				return err
			}

			go handler(m.Value)
		}
	}
}

func (c *ConsumerKafkaImpl) Close() {
	c.reader.Close()
}
