package repository

import (
	models "chi-demo/my_models"
	"context"
)

func (i ProductRepositoryImpl) Delete(ctx context.Context, id int64) error {
	product, err := models.FindProduct(ctx, i.db, id)
	if err != nil {
		return err
	}
	if _, err := product.Delete(ctx, i.db); err != nil {
		return err
	}
	return nil
}
