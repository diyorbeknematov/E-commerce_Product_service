package main

import (
	"log"
	"net"
	"product-service/config"
	"product-service/generated/mainservice"
	"product-service/logs"
	consumer "product-service/queue/kafka/concumer"
	"product-service/service"
	"product-service/storage"
	"product-service/storage/mongodb"
	"product-service/storage/redis"

	"google.golang.org/grpc"
)

func main() {
	logger := logs.InitLogger()
	logger.Info("Service started")

	db, err := mongodb.ConnectMongoDB()
	if err != nil {
		logger.Error("error in connecting to mongodb", "error", err)
		return
	}
	rdb := redis.RedisClient()

	cfg := config.Load()
	productConn, err := net.Listen("tcp", cfg.GRPC_PORT)
	if err != nil {
		logger.Error("error in listening to gRPC port", "error", err)
		return
	}
	storage := storage.NewStorage(rdb, db)
	s := grpc.NewServer()
	mainservice.RegisterMainServiceServer(s, service.NewProductService(storage, logger))

	go func() {
		log.Println("Starting kafka consumer...  Localhost:9092,")
		logger.Info("Starting kafka consumer...  Localhost:9092, Topic: order-created, Group")

		reader := consumer.NewConsumerKafka([]string{"localhost:9092"}, "order-created", "product-service", logger)
		defer reader.Close()
		serve := service.NewKafkaService(&reader, storage, logger)

		reader.ConsumeMessages(serve.CreateOrders)
	}()

	logger.Info("gRPC server started on port", "port", cfg.GRPC_PORT)
	log.Println("gRPC server started on port", "port", cfg.GRPC_PORT)
	if err := s.Serve(productConn); err != nil {
		logger.Error("gRPC server error", "error", err)
	}
}
