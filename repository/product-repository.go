package repository

import (
	"chi-demo/model"
	"context"
	"database/sql"
	"fmt"

	"github.com/sony/sonyflake"
)

type ProductRepositoryImpl struct {
	db    ContextExecutor
	idsnf *sonyflake.Sonyflake
}

type ProductRepository interface {
	GetOne(ctx context.Context, id int64) (model.Product, error)
	GetAll(ctx context.Context) ([]model.Product, error)
	Create(ctx context.Context, product model.Product) error
	Delete(ctx context.Context, id int64) error
}

// Executor can perform SQL queries.
type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// ContextExecutor can perform SQL queries with context
type ContextExecutor interface {
	Executor

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func New(db ContextExecutor) ProductRepository {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	if flake == nil {
		fmt.Printf("Couldn't generate sonyflake.NewSonyflake. Doesn't work on Go Playground due to fake time.\n")
	}

	return ProductRepositoryImpl{
		db:    db,
		idsnf: flake,
	}
}
