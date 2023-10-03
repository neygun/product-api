package repository

import (
	"chi-demo/model"
	models "chi-demo/my_models"
	"context"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (i ProductRepositoryImpl) Create(ctx context.Context, product model.Product) error {
	newID, err := i.idsnf.NextID()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	p := models.Product{
		ID:    int64(newID),
		Name:  product.Name,
		Price: product.Price,
	}

	if err := p.Insert(ctx, i.db, boil.Infer()); err != nil {
		return err
	}

	return nil
}
