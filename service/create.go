package service

import (
	"chi-demo/model"
	"context"
)

func (productServiceImpl ProductServiceImpl) Create(ctx context.Context, product model.Product) error {
	return productServiceImpl.productRepository.Create(ctx, product)
}
