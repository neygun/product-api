package repository

import (
	"chi-demo/model"
	models "chi-demo/my_models"
	"context"
)

func (i ProductRepositoryImpl) GetAll(ctx context.Context) ([]model.Product, error) {
	products, err := models.Products().All(ctx, i.db)
	if err != nil {
		return nil, err
	}

	result := make([]model.Product, len(products))

	for i, v := range products {
		result[i] = model.Product{
			ID:        v.ID,
			Name:      v.Name,
			Price:     v.Price,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
	}

	return result, nil
}
