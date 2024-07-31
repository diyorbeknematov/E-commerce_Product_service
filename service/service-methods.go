package service

import (
	"context"
	"fmt"
	"product-service/generated/categories"
	pb "product-service/generated/products"
	"product-service/generated/reviews"
)

// --------------------- productService ---------------------
func (p *productImpl) CreateProduct(ctx context.Context, in *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	res, err := p.storage.ProductRepository().CreateProduct(ctx, in)
	if err != nil {
		p.logger.Error("error in create product", "error", err)
		return nil, err
	}
	return res, nil
}

func (p *productImpl) UpdateProduct(ctx context.Context, in *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	res, err := p.storage.ProductRepository().UpdateProduct(ctx, in)
	if err != nil {
		p.logger.Error("error in update product", "error", err)
		return nil, err
	}
	return res, nil
}

func (p *productImpl) DeleteProduct(ctx context.Context, in *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	res, err := p.storage.ProductRepository().DeleteProduct(ctx, in)
	if err != nil {
		p.logger.Error("error in delete product", "error", err)
		return nil, err
	}
	return res, nil
}

func (p *productImpl) GetAllProduct(ctx context.Context, in *pb.GetAllProductRequest) (*pb.GetAllProductResponse, error) {
	resp, err := p.storage.ProductRepository().GetAllProducts(ctx, in)

	if err != nil {
		p.logger.Error("error in get all products", "error", err)
		return nil, err
	}

	return resp, nil
}
func (p *productImpl) GetByIdProduct(ctx context.Context, in *pb.GetByIdProductRequest) (*pb.GetByIdProductResponse, error) {
	result, err := p.storage.ProductRepository().GetProductByID(ctx, in.Id)
	if err != nil {
		p.logger.Error("error in get product by id", "error", err)
		return nil, err
	}
	return result, nil
}

func (p *productImpl) GetUserRecommendation(ctx context.Context, in *pb.Void) (*pb.GetRecommendationsResponse, error) {
	fmt.Println("Hello world")
	result, err := p.storage.ProductRepository().UserRecomendation(ctx)
	if err != nil {
		p.logger.Error("error in get user recommendation", "error", err)
		return nil, err
	}
	return result, nil
}

func (p *productImpl) GetPurchasedProducts(ctx context.Context, in *pb.GetPurchasedPRequest) (*pb.GetPurchasedPResponse, error) {
	result, err := p.storage.ProductRepository().GetUserBoughtProducts(ctx, in)
	if err != nil {
		p.logger.Error("error in get user bought products", "error", err)
		return nil, err
	}
	return result, nil
}

// --------------------- order and basket service ---------------------
func (p *productImpl) CreateOrder(ctx context.Context, in *pb.Order) (*pb.OrderResponse, error) {
	result, err := p.storage.OrderRepository().CreateOrder(ctx, in)
	if err != nil {
		p.logger.Error("error in create order", "error", err)
		return nil, err
	}
	return result, nil
}

func (p *productImpl) GetOrderByPId(ctx context.Context, in *pb.GetOrderByPIdRequest) (*pb.GetOrderByPIdResponse, error) {
	result, err := p.storage.OrderRepository().GetByProductId(ctx, in)
	if err != nil {
		p.logger.Error("error in get order by product id", "error", err)
		return nil, err
	}
	return result, nil
}

func (p *productImpl) AddToBasket(ctx context.Context, in *pb.AddToBasketRequest) (*pb.AddToBasketResponse, error) {
	result, err := p.storage.BasketRepository().AddToBasket(ctx, in)
	if err != nil {
		p.logger.Error("error in add to basket", "error", err)
		return nil, err
	}
	return result, nil
}

func (p *productImpl) GetBasketProducts(ctx context.Context, in *pb.GetBasketRequest) (*pb.BasketResponse, error) {
	products, err := p.storage.BasketRepository().GetFromBasket(ctx, in)
	if err != nil {
		p.logger.Error("error in get from basket", "error", err)
		return nil, err
	}
	result := &pb.BasketResponse{
		UserId:   in.UserId,
		Products: products,
	}

	return result, nil
}

func (p *productImpl) DeleteBasketProduct(ctx context.Context, in *pb.DeleteBasketRequest) (*pb.DeleteBasketResponse, error) {
	result, err := p.storage.BasketRepository().DeleteFromBasket(ctx, in)
	if err != nil {
		p.logger.Error("error in delete from basket", "error", err)
		return nil, err
	}
	return result, nil
}

// --------------------- categoryService ---------------------

func (p *productImpl) GetAllCategories(ctx context.Context, in *categories.GetAllCategoryRequest) (*categories.GetAllCategoryResponse, error) {
	result, err := p.storage.CategoryRepository().GetAllCategories(ctx, in)
	if err != nil {
		p.logger.Error("error in get all categories", "error", err)
		return nil, err
	}
	return result, nil
}

func (p *productImpl) CreateCategory(ctx context.Context, in *categories.CreateCategoryRequest) (*categories.CreateCategoryResponse, error) {
	result, err := p.storage.CategoryRepository().CreateCategory(ctx, in)
	if err != nil {
		p.logger.Error("error in create category", "error", err)
		return nil, err
	}
	return result, nil
}

func (p *productImpl) UpdateCategory(ctx context.Context, in *categories.UpdateCategoryRequest) (*categories.UpdateCategoryResponse, error) {
	result, err := p.storage.CategoryRepository().UpdateCategory(ctx, in)
	if err != nil {
		p.logger.Error("error in update category", "error", err)
		return nil, err
	}
	return result, nil
}

func (p *productImpl) DeleteCategory(ctx context.Context, in *categories.DeleteCategoryRequest) (*categories.DeleteCategoryResponse, error) {
	result, err := p.storage.CategoryRepository().DeleteCategory(ctx, in.Id)
	if err != nil {
		p.logger.Error("error in delete category", "error", err)
		return nil, err
	}
	return result, nil
}

// --------------------- reviewService ---------------------
func (p *productImpl) GetAllReviews(ctx context.Context, in *reviews.GetAllReviewsRequest) (*reviews.GetAllReviewsResponse, error) {
	result, err := p.storage.ReviewRepository().GetAllReviews(ctx, in)
	if err != nil {
		p.logger.Error("error in get all reviews", "error", err)
		return nil, err
	}
	return result, nil
}

func (p *productImpl) GetReviewsByProductId(ctx context.Context, in *reviews.GetReviewsByPIdRequest) (*reviews.GetReviewsByPIdResponse, error) {
	result, err := p.storage.ReviewRepository().GetAllReviews(ctx, &reviews.GetAllReviewsRequest{
		Offset:   in.GetOffset(),
		Limit:    in.GetLimit(),
		SearchBy: in.GetProductId(),
	})
	if err != nil {
		p.logger.Error("error in get reviews by product id", "error", err)
		return nil, err
	}
	return &reviews.GetReviewsByPIdResponse{
		Reviews: result.Reviews,
		Total:   result.Total,
		Page:    result.Page,
		Limit:   result.Limit,
	}, nil
}

func (p *productImpl) CreateReview(ctx context.Context, in *reviews.CreateReviewRequest) (*reviews.CreateReviewResponse, error) {
	result, err := p.storage.ReviewRepository().CreateReview(ctx, in)

	if err != nil {
		p.logger.Error("error in create review", "error", err)
		return nil, err
	}
	return result, nil
}

func (p *productImpl) UpdateReview(ctx context.Context, in *reviews.UpdateReviewRequest) (*reviews.UpdateReviewResponse, error) {
	result, err := p.storage.ReviewRepository().UpdateReview(ctx, in)
	if err != nil {
		p.logger.Error("error in update review", "error", err)
		return nil, err
	}
	return result, nil
}

func (p *productImpl) DeleteReview(ctx context.Context, in *reviews.DeleteReviewRequest) (*reviews.DeleteReviewResponse, error) {
	result, err := p.storage.ReviewRepository().DeleteReview(ctx, in)
	if err != nil {
		p.logger.Error("error in delete review", "error", err)
		return nil, err
	}
	return result, nil
}
