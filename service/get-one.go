package service

import (
	"chi-demo/model"
	"context"
)

func (productServiceImpl ProductServiceImpl) GetOne(ctx context.Context, id int64) (model.Product, error) {
	return productServiceImpl.ProductRepository.GetOne(ctx, id)
}
