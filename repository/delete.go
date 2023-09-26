package repository

import (
	models "chi-demo/my_models"
	"context"
)

func (i ProductRepositoryImpl) Delete(ctx context.Context, id int64) error {
	product, _ := models.FindProduct(ctx, i.db, id)
	_, err := product.Delete(ctx, i.db)
	return err
}
