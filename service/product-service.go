package service

import (
	"chi-demo/model"
	"chi-demo/repository"
	"context"
)

type ProductService interface {
	GetOne(ctx context.Context, id int64) (model.Product, error)
	GetAll(ctx context.Context) ([]model.Product, error)
	Create(ctx context.Context, product model.Product) error
	Delete(ctx context.Context, id int64) error
}

type ProductServiceImpl struct {
	productRepository repository.ProductRepository
}

func New(productRepository repository.ProductRepository) ProductService {
	return ProductServiceImpl{
		productRepository: productRepository,
	}
}
