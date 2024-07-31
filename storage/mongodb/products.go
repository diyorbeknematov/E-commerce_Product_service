package mongodb

import (
	"context"
	"fmt"
	pb "product-service/generated/products"

	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *pb.CreateProductRequest) (*pb.CreateProductResponse, error)
	GetProductByID(ctx context.Context, id string) (*pb.GetByIdProductResponse, error)
	UpdateProduct(ctx context.Context, product *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error)
	DeleteProduct(ctx context.Context, in *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error)
	UserRecomendation(ctx context.Context) (*pb.GetRecommendationsResponse, error)
	GetUserBoughtProducts(ctx context.Context, in *pb.GetPurchasedPRequest) (*pb.GetPurchasedPResponse, error)
	GetAllProducts(ctx context.Context, fProduct *pb.GetAllProductRequest) (*pb.GetAllProductResponse, error)
}

type productRepositoryImpl struct {
	coll *mongo.Collection
}

func NewProductRepository(db *mongo.Database) ProductRepository {
	return &productRepositoryImpl{coll: db.Collection("product")}
}

func (p *productRepositoryImpl) CreateProduct(ctx context.Context, product *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	_, err := p.coll.InsertOne(ctx, bson.M{
		"_id":         uuid.NewString(),
		"category_id": product.CategoryId,
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"stock":       product.Stock,
		"images":      bson.A{product.Images},
		"discount": bson.M{
			"status":         product.Discount.Status,
			"discount_price": product.Discount.DiscountPrice,
		},
		"created_at": time.Now().Format("2006-01-02 15:04:05"),
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
		"deleted_at": "",
	})

	return &pb.CreateProductResponse{
		Success: true,
		Message: "Product created successfully",
	}, err
}

func (p *productRepositoryImpl) GetProductByID(ctx context.Context, id string) (*pb.GetByIdProductResponse, error) {
	var product pb.GetByIdProductResponse
	err := p.coll.FindOne(ctx, bson.M{"$and": []bson.M{{"_id": id}, {"deleted_at": ""}}}).Decode(&product)

	if err == mongo.ErrNilDocument {
		return nil, err
	}

	return &product, err
}

func (p *productRepositoryImpl) UpdateProduct(ctx context.Context, product *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	filter := bson.M{
		"$and": []bson.M{
			{"_id": product.Id},
			{"deleted_at": ""},
		}}

	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
			"discount": bson.M{
				"price":  product.Discount.DiscountPrice,
				"status": product.Discount.Status,
			},
			"updated_at": time.Now().Format("2006-01-02 15:04:05"),
		},
		"$push": product.Images,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedProduct pb.UpdateProductResponse
	err := p.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedProduct)
	if err != nil {
		return nil, err
	}

	return &updatedProduct, nil
}

func (p *productRepositoryImpl) DeleteProduct(ctx context.Context, in *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	_, err := p.coll.UpdateByID(ctx,
		bson.M{"&and": []bson.M{
			{"_id": in.Id},
			{"deleted_at": ""},
		}},
		bson.M{"$set": bson.M{
			"deleted_at": time.Now().Format("2006-01-02 15:04:05"),
		}})

	if err != nil {
		return nil, err
	}

	return &pb.DeleteProductResponse{
		Success: true,
		Message: "Product deleted successfully",
	}, nil
}

func (p *productRepositoryImpl) UserRecomendation(ctx context.Context) (*pb.GetRecommendationsResponse, error) {
	fmt.Println("Hello world")
	filter := bson.D{
		{Key: "discount.status", Value: true},
		{Key: "deleted_at", Value: ""},
	}

	cursor, err := p.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var products []*pb.Product
	for cursor.Next(ctx) {
		var product pb.Product
		err := cursor.Decode(&product)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		products = append(products, &product)
	}

	return &pb.GetRecommendationsResponse{
		Products: products,
	}, nil
}

func (p *productRepositoryImpl) GetUserBoughtProducts(ctx context.Context, in *pb.GetPurchasedPRequest) (*pb.GetPurchasedPResponse, error) {
	filter := bson.D{
		{Key: "deleted_at", Value: ""},
		{Key: "user_id", Value: in.GetUserId()},
	}

	total, err := p.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	opts := options.Find()
	opts.SetLimit(in.GetLimit())
	opts.SetSkip(in.GetPage())

	cursor, err := p.coll.Find(ctx, filter, opts)
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

	return &pb.GetPurchasedPResponse{
		Orders: orders,
		Total:  total,
		Limit:  in.GetLimit(),
		Page:   in.GetPage(),
	}, nil
}

func (p *productRepositoryImpl) GetAllProducts(ctx context.Context, fProduct *pb.GetAllProductRequest) (*pb.GetAllProductResponse, error) {
	// Starting the pipeline with lookup and unwind stages
	var pipeline = mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "categories"},
			{Key: "localField", Value: "category_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "category"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$category"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "deleted_at", Value: ""}}}},
	}

	// Adding filters and sorts
	filters, sorts := createFilters(fProduct)

	for _, filter := range filters {
		pipeline = append(pipeline, filter)
	}

	for _, sort := range sorts {
		pipeline = append(pipeline, sort)
	}

	// Counting total results
	countPipeline := append(pipeline, bson.D{{Key: "$count", Value: "total"}})

	countCursor, err := p.coll.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	var countResult []bson.M
	if err := countCursor.All(ctx, &countResult); err != nil {
		return nil, err
	}

	var total int64
	if len(countResult) > 0 {
		totalInt32 := countResult[0]["total"].(int32)
		total = int64(totalInt32)
	}

	// Adding pagination
	pipeline = append(pipeline, bson.D{{Key: "$skip", Value: (fProduct.GetPage() - 1) * fProduct.GetLimit()}})
	pipeline = append(pipeline, bson.D{{Key: "$limit", Value: fProduct.GetLimit()}})

	cursor, err := p.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*pb.Product
	for cursor.Next(ctx) {
		var product pb.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return &pb.GetAllProductResponse{
		Products: products,
		Total:    total,
		Limit:    fProduct.Limit,
		Offset:   fProduct.Page,
	}, nil
}

func createFilters(in *pb.GetAllProductRequest) ([]bson.D, []bson.D) {
	var filters []bson.D
	var sorts []bson.D

	if in.GetName() != "" {
		filters = append(filters, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "name", Value: bson.D{
					{Key: "$regex", Value: ".*" + in.GetName() + ".*"},
					{Key: "$options", Value: "i"},
				}},
			}},
		})
	}
	if in.GetCategory() != "" {
		filters = append(filters, bson.D{{Key: "$match", Value: bson.D{{Key: "category.name", Value: in.GetCategory()}}}})
	}
	if in.GetDiscount() {
		filters = append(filters, bson.D{{Key: "$match", Value: bson.D{{Key: "discount.status", Value: in.Discount}}}})
	}
	if in.GetNewest() {
		sorts = append(sorts, bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}})
	}
	if in.GetPriceOrder() != 0 {
		sorts = append(sorts, bson.D{{Key: "$sort", Value: bson.D{{Key: "price", Value: in.GetPriceOrder()}}}})
	}

	if in.GetRatingOrder() != 0 {
		sorts = append(sorts, bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "comments"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "product_id"},
			{Key: "as", Value: "comments"},
		}}})
		sorts = append(sorts, bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$comments"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}})
		sorts = append(sorts, bson.D{{Key: "$sort", Value: bson.D{{Key: "comments.rating", Value: in.GetRatingOrder()}}}})
	}

	if in.GetCommentOrder() != 0 {
		filters = append(filters, bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "comments"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "product_id"},
			{Key: "as", Value: "comments"},
		}}})
		filters = append(filters, bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "comment_count", Value: bson.D{{Key: "$size", Value: "$comments"}}},
		}}})
		sorts = append(sorts, bson.D{{Key: "$sort", Value: bson.D{{Key: "comment_count", Value: in.GetCommentOrder()}}}})
		sorts = append(sorts, bson.D{{Key: "$sort", Value: bson.D{{Key: "rating", Value: -1}}}})
	}

	return filters, sorts
}
