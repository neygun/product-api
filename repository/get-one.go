package repository

import (
	"chi-demo/model"
	models "chi-demo/my_models"
	"context"
	"log"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (i ProductRepositoryImpl) GetOne(ctx context.Context, id int64) (model.Product, error) {
	product, err := models.Products(qm.Where("id=?", id)).One(ctx, i.db)
	if err != nil {
		log.Println(err)
		return model.Product{}, err
	}
	return model.Product{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}, nil
}
