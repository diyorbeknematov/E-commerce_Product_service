package storage

import (
	"product-service/storage/mongodb"
	rdb "product-service/storage/redis"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type IStorage interface {
	BasketRepository() rdb.BasketRepository
	ProductRepository() mongodb.ProductRepository
	CategoryRepository() mongodb.CategoryRepository
	ReviewRepository() mongodb.ReviewRepository
	OrderRepository() mongodb.OrderRepository
}

type storageImpl struct {
	rdb *redis.Client
	db  *mongo.Database
}

func NewStorage(rdb *redis.Client, db *mongo.Database) IStorage {
	return &storageImpl{
		rdb: rdb,
		db:  db,
	}
}

func (s *storageImpl) BasketRepository() rdb.BasketRepository {
	return rdb.NewBasketRepository(s.rdb)
}

func (s *storageImpl) ProductRepository() mongodb.ProductRepository {
	return mongodb.NewProductRepository(s.db)
}

func (s *storageImpl) CategoryRepository() mongodb.CategoryRepository {
	return mongodb.NewCategoryRepository(s.db)
}

func (s *storageImpl) ReviewRepository() mongodb.ReviewRepository {
	return mongodb.NewReviewRepository(s.db)
}

func (s *storageImpl) OrderRepository() mongodb.OrderRepository {
	return mongodb.NewBoughtProductCollection(s.db)
}