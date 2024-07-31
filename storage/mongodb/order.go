package mongodb

import (
	"context"
	pb "product-service/generated/products"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *pb.Order) (*pb.OrderResponse, error)
	GetByProductId(ctx context.Context, in *pb.GetOrderByPIdRequest) (*pb.GetOrderByPIdResponse, error)
}

type boughtProductImpl struct {
	coll *mongo.Collection
}

func NewBoughtProductCollection(db *mongo.Database) OrderRepository {
	return &boughtProductImpl{
		coll: db.Collection("bought_products"),
	}
}

func (b *boughtProductImpl) CreateOrder(ctx context.Context, order *pb.Order) (*pb.OrderResponse, error) {
	_, err := b.coll.InsertOne(ctx, bson.D{
		{Key: "_id", Value: order.GetId()},
		{Key: "user_id", Value: order.GetUserId()},
		{Key: "purchase_date", Value: order.GetPurchaseDate()},
		{Key: "quantity", Value: order.GetQuantity()},
		{Key: "price", Value: order.GetPrice()},
	})
	return &pb.OrderResponse{
		ProductId: order.GetId(),
	}, err
}

func (b *boughtProductImpl) GetByProductId(ctx context.Context, in *pb.GetOrderByPIdRequest) (*pb.GetOrderByPIdResponse, error) {
	filter := bson.D{
		{Key: "product_id", Value: in.GetProductId()},
	}

	total, err := b.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	opts := options.Find()
	opts.SetLimit(in.GetLimit())
	opts.SetSkip(in.GetPage())
	cursor, err := b.coll.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var orders []*pb.Order

	for cursor.Next(ctx) {
		var order pb.Order
		err := cursor.Decode(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return &pb.GetOrderByPIdResponse{
		Orders: orders,
		Total:  total,
		Limit:  in.GetLimit(),
		Offset: in.GetPage(),
	}, nil
}
