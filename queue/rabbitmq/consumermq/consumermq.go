package consumermq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerMQ interface {
	// ConsumeMessages(queue string) ([]string, error)
	// Close() error
}

type ConsumerMQImpl struct {
	consumer *amqp.Channel
}

func NewConsumerMQ(url string) (ConsumerMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &ConsumerMQImpl{consumer: ch}, nil
}

// func (c *ConsumerMQImpl) ConsumeMessages(queue string) ([]string, error) {

// }

func (c *ConsumerMQImpl) Close() error {
	return c.consumer.Close()
}
