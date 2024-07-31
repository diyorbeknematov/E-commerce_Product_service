package service

import (
	"context"
	"log/slog"
	"product-service/generated/products"
	consumer "product-service/queue/kafka/concumer"
	"product-service/storage"

	"google.golang.org/protobuf/proto"
)

type KafkaService interface {
	CreateOrders(messageset []byte)
}

type KafkaServiceImpl struct {
	consumer *consumer.KafkaConsumer
	storage  storage.IStorage
	logger   *slog.Logger
}

func NewKafkaService(consumer *consumer.KafkaConsumer, storage storage.IStorage, logger *slog.Logger) KafkaService {
	return &KafkaServiceImpl{
		consumer: consumer,
		storage:  storage,
		logger:   logger,
	}
}

func (s *KafkaServiceImpl) CreateOrders(messageset []byte) {
	var order products.OrderRequest
	err := proto.Unmarshal(messageset, &order)
	if err != nil {
		s.logger.Error("error in unmarshalling message", "error", err)
		return
	}

	product, err := s.storage.BasketRepository().GetFromBasketById(context.Background(), order.UserId, order.ProductId)
	if err != nil {
		s.logger.Error("error in get from basket", "error", err)
		return
	}

	resp, err := s.storage.OrderRepository().CreateOrder(context.Background(), &products.Order{
		Id:           product.Id,
		UserId:       product.UserId,
		PurchaseDate: product.PurchaseDate,
		Quantity:     product.Quantity,
		Price:        product.Price,
	})
	if err != nil {
		s.logger.Error("error in create order", "error", err)
		return
	}

	s.logger.Info("order created successfully", "order", resp)
}
