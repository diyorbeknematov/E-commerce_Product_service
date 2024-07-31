package service

import (
	"context"
	"log/slog"
	"product-service/generated/categories"
	"product-service/generated/mainservice"
	"product-service/generated/products"
	"product-service/generated/reviews"
	"product-service/storage"
)

type ProductService interface {
	// Mahsulotlar bilan bog'liq endpointlar
	CreateProduct(context.Context, *products.CreateProductRequest) (*products.CreateProductResponse, error)
	UpdateProduct(context.Context, *products.UpdateProductRequest) (*products.UpdateProductResponse, error)
	DeleteProduct(context.Context, *products.DeleteProductRequest) (*products.DeleteProductResponse, error)
	GetAllProduct(context.Context, *products.GetAllProductRequest) (*products.GetAllProductResponse, error)
	GetByIdProduct(context.Context, *products.GetByIdProductRequest) (*products.GetByIdProductResponse, error)
	CreateOrder(context.Context, *products.Order) (*products.OrderResponse, error)
	GetOrderByPId(context.Context, *products.GetOrderByPIdRequest) (*products.GetOrderByPIdResponse, error)
	AddToBasket(context.Context, *products.AddToBasketRequest) (*products.AddToBasketResponse, error)
	GetBasketProducts(context.Context, *products.GetBasketRequest) (*products.BasketResponse, error)
	DeleteBasketProduct(context.Context, *products.DeleteBasketRequest) (*products.DeleteBasketResponse, error)
	GetUserRecommendation(context.Context, *products.Void) (*products.GetRecommendationsResponse, error)
	GetPurchasedProducts(context.Context, *products.GetPurchasedPRequest) (*products.GetPurchasedPResponse, error)
	// Kategoriyalar bilan bog'liq endpointlar
	GetAllCategories(context.Context, *categories.GetAllCategoryRequest) (*categories.GetAllCategoryResponse, error)
	CreateCategory(context.Context, *categories.CreateCategoryRequest) (*categories.CreateCategoryResponse, error)
	UpdateCategory(context.Context, *categories.UpdateCategoryRequest) (*categories.UpdateCategoryResponse, error)
	DeleteCategory(context.Context, *categories.DeleteCategoryRequest) (*categories.DeleteCategoryResponse, error)
	// Sharhlar bilan bog'liq endpointlar
	GetAllReviews(context.Context, *reviews.GetAllReviewsRequest) (*reviews.GetAllReviewsResponse, error)
	GetReviewsByProductId(context.Context, *reviews.GetReviewsByPIdRequest) (*reviews.GetReviewsByPIdResponse, error)
	CreateReview(context.Context, *reviews.CreateReviewRequest) (*reviews.CreateReviewResponse, error)
	UpdateReview(context.Context, *reviews.UpdateReviewRequest) (*reviews.UpdateReviewResponse, error)
	DeleteReview(context.Context, *reviews.DeleteReviewRequest) (*reviews.DeleteReviewResponse, error)
}

type productImpl struct {
	mainservice.UnimplementedMainServiceServer
	logger  *slog.Logger
	storage storage.IStorage
}

func NewProductService(storage storage.IStorage, logger *slog.Logger) *productImpl {
	return &productImpl{
		storage: storage,
		logger: logger,
	}
}
