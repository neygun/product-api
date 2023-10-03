package repository

import (
	"chi-demo/db"
	"chi-demo/model"
	"context"
	"fmt"

	"github.com/sony/sonyflake"
)

type ProductRepositoryImpl struct {
	db    db.ContextExecutor
	idsnf *sonyflake.Sonyflake
}

type ProductRepository interface {
	GetOne(ctx context.Context, id int64) (model.Product, error)
	GetAll(ctx context.Context) ([]model.Product, error)
	Create(ctx context.Context, product model.Product) error
	Delete(ctx context.Context, id int64) error
}

func New(db db.ContextExecutor) ProductRepository {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	if flake == nil {
		fmt.Printf("Couldn't generate sonyflake.NewSonyflake. Doesn't work on Go Playground due to fake time.\n")
	}

	return ProductRepositoryImpl{
		db:    db,
		idsnf: flake,
	}
}
