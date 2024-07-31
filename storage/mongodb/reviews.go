package mongodb

import (
	"context"
	pb "product-service/generated/reviews"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReviewRepository interface {
	GetAllReviews(ctx context.Context, review *pb.GetAllReviewsRequest) (*pb.GetAllReviewsResponse, error)
	CreateReview(ctx context.Context, review *pb.CreateReviewRequest) (*pb.CreateReviewResponse, error)
	UpdateReview(ctx context.Context, review *pb.UpdateReviewRequest) (*pb.UpdateReviewResponse, error)
	DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.DeleteReviewResponse, error)
	GetReviewByID(ctx context.Context, id string) (*pb.GetReviewsByPIdResponse, error)
}

type reviewRepositoryImpl struct {
	coll *mongo.Collection
}

func NewReviewRepository(db *mongo.Database) ReviewRepository {
	return &reviewRepositoryImpl{coll: db.Collection("comments")}
}

func (r *reviewRepositoryImpl) GetAllReviews(ctx context.Context, review *pb.GetAllReviewsRequest) (*pb.GetAllReviewsResponse, error) {
	filter := bson.D{
		{Key: "deleted_at", Value: ""},
	}

	if review.SearchBy != "" {
		filter = append(filter, bson.E{Key: "$or", Value: bson.A{
			bson.D{{Key: "product_id", Value: review.SearchBy}},
			bson.D{{Key: "user_id", Value: review.SearchBy}},
		}})
	}

	// Umumiy sonini olish
	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	opts := options.Find()
	opts.SetSkip(review.Offset)

	opts.SetLimit(review.Limit)

	if review.SortBy != 0 {
		opts.SetSort(bson.D{{Key: "created_at", Value: review.SortBy}})
	}

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []*pb.Review
	for cursor.Next(ctx) {
		var r pb.Review
		err := cursor.Decode(&r)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &r)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &pb.GetAllReviewsResponse{
		Reviews: reviews,
		Total:   total,
		Page:    review.Offset,
		Limit:   review.Limit,
	}, nil
}

func (r *reviewRepositoryImpl) GetReviewByID(ctx context.Context, id string) (*pb.GetReviewsByPIdResponse, error) {
	filter := bson.M{
		"$and": bson.D{
			{Key: "deleted_at", Value: ""},
			{Key: "_id", Value: id},
		},
	}
	var review pb.GetReviewsByPIdResponse
	err := r.coll.FindOne(ctx, filter).Decode(&review)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepositoryImpl) CreateReview(ctx context.Context, review *pb.CreateReviewRequest) (*pb.CreateReviewResponse, error) {
	res, err := r.coll.InsertOne(ctx, bson.D{
		{Key: "_id", Value: uuid.NewString()},
		{Key: "product_id", Value: review.ProductId},
		{Key: "user_id", Value: review.UserId},
		{Key: "rating", Value: review.Rating},
		{Key: "comment", Value: review.Comment},
		{Key: "created_at", Value: time.Now().Format("2006-01-02 15:04:05")},
		{Key: "updated_at", Value: time.Now().Format("2006-01-02 15:04:05")},
		{Key: "deleted_at", Value: ""},
	})

	if err != nil {
		return nil, err
	}

	return &pb.CreateReviewResponse{
		Id:        res.InsertedID.(string),
		UserId:    review.UserId,
		ProductId: review.ProductId,
		Rating:    review.Rating,
		Comment:   review.Comment,
	}, nil
}

func (r *reviewRepositoryImpl) UpdateReview(ctx context.Context, review *pb.UpdateReviewRequest) (*pb.UpdateReviewResponse, error) {
	filter := bson.M{
		"$and": bson.D{
			{Key: "_id", Value: review.Id},
			{Key: "user_id", Value: review.UserId},
			{Key: "product_id", Value: review.ProductId},
			{Key: "deleted_at", Value: ""},
		},
	}

	update := bson.M{
		"$set": bson.D{
			{Key: "rating", Value: review.Rating},
			{Key: "comment", Value: review.Comment},
			{Key: "updated_at", Value: time.Now().Format("2006-01-02 15:04:05")},
		},
	}
	var updatedReview pb.UpdateReviewResponse
	if err := r.coll.FindOneAndUpdate(ctx, filter, update).Decode(&updatedReview); err != nil {
		return nil, err
	}

	return &updatedReview, nil
}

func (r *reviewRepositoryImpl) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.DeleteReviewResponse, error) {
	filter := bson.M{
		"$and": bson.D{
			{Key: "_id", Value: req.Id},
			{Key: "user_id", Value: req.UserId},
			{Key: "product_id", Value: req.ProductId},
			{Key: "deleted_at", Value: ""},
		},
	}

	update := bson.M{
		"$set": bson.D{
			{Key: "deleted_at", Value: time.Now().Format("2006-01-02 15:04:05")},
		},
	}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.ModifiedCount == 0 {
		return nil, nil
	}

	return &pb.DeleteReviewResponse{
		Success: true,
		Message: "Review deleted successfully",
	}, nil
}
