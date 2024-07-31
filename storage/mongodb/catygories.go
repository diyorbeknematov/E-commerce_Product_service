package mongodb

import (
	"context"
	pb "product-service/generated/categories"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryRepository interface {
	GetAllCategories(ctx context.Context, category *pb.GetAllCategoryRequest) (*pb.GetAllCategoryResponse, error)
	CreateCategory(ctx context.Context, category *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error)
	UpdateCategory(ctx context.Context, category *pb.UpdateCategoryRequest) (*pb.UpdateCategoryResponse, error)
	DeleteCategory(ctx context.Context, id string) (*pb.DeleteCategoryResponse, error)
}

type categoryRepositoryImpl struct {
	coll *mongo.Collection
}

func NewCategoryRepository(db *mongo.Database) CategoryRepository {
	return &categoryRepositoryImpl{coll: db.Collection("categories")}
}

func (c *categoryRepositoryImpl) GetAllCategories(ctx context.Context, category *pb.GetAllCategoryRequest) (*pb.GetAllCategoryResponse, error) {
	var filter = bson.M{
		"$and": []bson.M{
			{"deleted_at": ""},
		},
	}
	total, err := c.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	opts := options.Find()
	opts.SetSkip(category.Offset)
	opts.SetLimit(category.Limit)

	var categories []*pb.Category
	cursor, err := c.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var category pb.Category
		err := cursor.Decode(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &pb.GetAllCategoryResponse{
		Catygories: categories,
		Offset: category.Offset,
		Limit:     category.Limit,
        Total:     total,
	}, nil
}

func (c *categoryRepositoryImpl) CreateCategory(ctx context.Context, category *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	res, err := c.coll.InsertOne(ctx, bson.M{
		"_id":         uuid.NewString(),
		"name":        category.Name,
		"description": category.Description,
		"created_at":  time.Now().Format("2006-01-02 15:04:05"),
		"updated_at":  time.Now().Format("2006-01-02 15:04:05"),
		"deleted_at":  "",
	})
	if err != nil {
		return &pb.CreateCategoryResponse{
			Status:  false,
			Message: "Failed to create category",
			Id:      "",
		}, err
	}

	id := res.InsertedID.(string)
	return &pb.CreateCategoryResponse{
		Status:  true,
		Message: "Category created successfully",
		Id: id,
	}, nil
}

func (c *categoryRepositoryImpl) UpdateCategory(ctx context.Context, category *pb.UpdateCategoryRequest) (*pb.UpdateCategoryResponse, error) {
	filter := bson.M{
		"_id":        category.Id,
		"deleted_at": "",
	}

	update := bson.M{
		"$set": []bson.M{
			{"name": category.Name},
			{"description": category.Description},
			{"updated_at": time.Now().Format("2006-01-02 15:04:05")},
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedCategory pb.UpdateCategoryResponse

	err := c.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(update)

	return &updatedCategory, err
}

func (c *categoryRepositoryImpl) DeleteCategory(ctx context.Context, id string) (*pb.DeleteCategoryResponse, error) {
	filter := bson.M{
		"_id":        id,
		"deleted_at": "",
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	result, err := c.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.ModifiedCount == 0 {
		return &pb.DeleteCategoryResponse{
			Status:  false,
			Message: "Category not found",
		}, nil
	}

	return &pb.DeleteCategoryResponse{
		Status:  true,
		Message: "Category deleted successfully",
	}, nil
}
