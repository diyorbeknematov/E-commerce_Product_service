package redis

import (
	"context"
	"encoding/json"
	"errors"
	"product-service/generated/products"

	"github.com/redis/go-redis/v9"
)

type BasketRepository interface {
	AddToBasket(ctx context.Context, in *products.AddToBasketRequest) (*products.AddToBasketResponse, error)
	GetFromBasket(ctx context.Context, in *products.GetBasketRequest) ([]*products.Order, error)
	GetFromBasketById(ctx context.Context, userId, productId string) (*products.Order, error)
	DeleteFromBasket(ctx context.Context, in *products.DeleteBasketRequest) (*products.DeleteBasketResponse, error)
}

type basketImpl struct {
	redis *redis.Client
}

func NewBasketRepository(redis *redis.Client) BasketRepository {
	return &basketImpl{redis: redis}
}

func (b *basketImpl) AddToBasket(ctx context.Context, in *products.AddToBasketRequest) (*products.AddToBasketResponse, error) {
	data, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	err = b.redis.RPush(ctx, in.UserId, data).Err()
	if err != nil {
		return nil, err
	}

	return &products.AddToBasketResponse{
		Status:  true,
		Message: "added to basket",
	}, nil
}

func (b *basketImpl) GetFromBasket(ctx context.Context, in *products.GetBasketRequest) ([]*products.Order, error) {

	list, err := b.redis.LRange(ctx, in.UserId, 0, -1).Result()

	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("not found")
	}

	var result []*products.Order 
	for _, v := range list {
		var order products.Order
        err := json.Unmarshal([]byte(v), &order)
        if err != nil {
            return nil, err
        }
        result = append(result, &order)
	}
	return result, nil
}

func (b *basketImpl) GetFromBasketById(ctx context.Context, userId, productId string) (*products.Order, error) {

	list, err := b.redis.LRange(ctx, userId, 0, -1).Result()

	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("not found")
	}

	var result products.Order

	for i := 0; i < len(list); i++ {
		if list[i] == productId {
			err := json.Unmarshal([]byte(list[i]), &result)
			if err != nil {
				return nil, err
			}
			return &result, nil
		}
	}

	return nil, errors.New("not found in basket")
}

func (b *basketImpl) DeleteFromBasket(ctx context.Context, in *products.DeleteBasketRequest) (*products.DeleteBasketResponse, error) {

	list, err := b.redis.LRange(ctx, in.UserId, 0, -1).Result()

	if err != nil {
		return nil, err
	}

	check := true
	for i := 0; i < len(list); i++ {
		var product products.Order
		err := json.Unmarshal([]byte(list[i]), &product)
		if err!= nil {
            return nil, err
        }
		if product.Id == in.ProductId {
			list = append(list[:i], list[i+1:]...)
			check = false
		}
	}
	if check {
		return &products.DeleteBasketResponse{
			Success: false,
			Message: "Product notfound",
		}, nil
	}

	err = b.redis.Del(ctx, in.UserId).Err()
	if err != nil {
		return nil, err
	}

	err = b.redis.RPush(ctx, in.UserId, list).Err()
	if err != nil {
		return nil, err
	}

	return &products.DeleteBasketResponse{
		Success: true,
		Message: "Product deleted",
	}, nil
}
