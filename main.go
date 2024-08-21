package main

import (
	"context"
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
		log.Println(err)
		logger.Error("error in connecting to mongodb", "error", err)
		return
	}
	rdb := redis.RedisClient()

	cfg := config.Load()
	productConn, err := net.Listen("tcp", cfg.GRPC_PORT)
	if err != nil {
		log.Println(err)
		logger.Error("error in listening to gRPC port", "error", err)
		return
	}
	storage := storage.NewStorage(rdb, db)
	s := grpc.NewServer()
	mainservice.RegisterMainServiceServer(s, service.NewProductService(storage, logger))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		logger.Info("Starting kafka consumer...  :9092, Topic: order-created, Group")

		reader := consumer.NewConsumerKafka([]string{"kafka:9092"}, "order-created", "product-service", logger)
		defer reader.Close()
		serve := service.NewKafkaService(&reader, storage, logger)

		err := reader.ConsumeMessages(ctx, serve.CreateOrders)
		if err != nil {
			log.Println("error in consuming kafka messages", "error", err)
			logger.Error("error in consuming kafka messages", "error", err)
            return
		}
	}()

	logger.Info("gRPC server started on port", "port", cfg.GRPC_PORT)
	log.Println("gRPC server started on port", "port", cfg.GRPC_PORT)
	if err := s.Serve(productConn); err != nil {
		logger.Error("gRPC server error", "error", err)
	}
}
