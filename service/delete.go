package service

import "context"

func (productServiceImpl ProductServiceImpl) Delete(ctx context.Context, id int64) error {
	return productServiceImpl.productRepository.Delete(ctx, id)
}
