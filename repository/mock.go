package repository

import (
	"chi-demo/model"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (r MockRepository) GetOne(ctx context.Context, id int64) (model.Product, error) {
	args := r.Called(ctx, id)

	return args.Get(0).(model.Product), args.Error(1)
}

func (r MockRepository) GetAll(ctx context.Context) ([]model.Product, error) {
	args := r.Called(ctx)
	return args.Get(0).([]model.Product), args.Error(1)

}

func (r MockRepository) Create(ctx context.Context, product model.Product) error {
	args := r.Called(ctx, product)
	return args.Error(0)
}

func (r MockRepository) Delete(ctx context.Context, id int64) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}
