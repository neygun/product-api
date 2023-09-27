package service

import (
	"chi-demo/model"
	"context"
)

func (productServiceImpl ProductServiceImpl) GetAll(ctx context.Context) ([]model.Product, error) {
	return productServiceImpl.productRepository.GetAll(ctx)
}
