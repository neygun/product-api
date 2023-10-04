package repository

import (
	"chi-demo/db"
	"chi-demo/model"
	"chi-demo/repository/testdata"
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestImpl_Create(t *testing.T) {
	type args struct {
		expDBFailed bool
		expErr      error
	}

	tcs := map[string]args{
		"success": {
			expErr: nil,
		},
		"error: db failed": {
			expDBFailed: true,
			expErr:      errors.New("models: unable to insert into product: sql: database is closed"),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			TestWithTxDB(t, func(tx db.ContextExecutor) {
				// Given
				repo := New(tx)
				if tc.expDBFailed {
					dbMock, _, _ := sqlmock.New()
					dbMock.Close()
					repo = New(dbMock)
				}
				testdata.LoadTestSQLFile(t, tx, "testdata/create_product.sql")

				// When
				product := model.Product{
					ID:    1,
					Name:  "test",
					Price: 1,
				}
				err := repo.Create(ctx, product)

				// Then
				if tc.expErr != nil {
					require.EqualError(t, err, tc.expErr.Error())
				} else {
					require.NoError(t, err)
				}
			})
		})
	}
}

func TestImpl_GetById(t *testing.T) {
	type args struct {
		givenID     int64
		expDBFailed bool
		expRs       model.Product
		expErr      error
	}

	tcs := map[string]args{
		"success": {
			givenID: 1,
			expRs: model.Product{
				ID:    1,
				Name:  "test",
				Price: 1,
			},
		},
		"error: not found": {
			givenID: 1000,
			expErr:  sql.ErrNoRows,
		},
		"error: db failed": {
			givenID:     1,
			expDBFailed: true,
			expErr:      errors.New("models: failed to execute a one query for product: bind failed to execute query: sql: database is closed"),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			TestWithTxDB(t, func(tx db.ContextExecutor) {
				// Given
				repo := New(tx)
				if tc.expDBFailed {
					dbMock, _, _ := sqlmock.New()
					dbMock.Close()
					repo = New(dbMock)
				}
				testdata.LoadTestSQLFile(t, tx, "testdata/get_product_by_id.sql")

				// When
				result, err := repo.GetOne(ctx, tc.givenID)

				// Then
				if tc.expErr != nil {
					require.EqualError(t, err, tc.expErr.Error())
				} else {
					require.NoError(t, err)
					require.NotZero(t, result.CreatedAt)
					require.NotZero(t, result.UpdatedAt)
					require.NotZero(t, result.DeletedAt)
					if !cmp.Equal(tc.expRs, result,
						cmpopts.IgnoreFields(model.Product{}, "CreatedAt", "UpdatedAt", "DeletedAt")) {
						t.Errorf("\n order mismatched. \n expected: %+v \n got: %+v \n diff: %+v", tc.expRs, result,
							cmp.Diff(tc.expRs, result, cmpopts.IgnoreFields(model.Product{}, "CreatedAt", "UpdatedAt", "DeletedAt")))
						t.FailNow()
					}
				}
			})
		})
	}
}

func TestImpl_GetAll(t *testing.T) {
	type args struct {
		isEmpty     bool
		expDBFailed bool
		expRs       []model.Product
		expErr      error
	}

	tcs := map[string]args{
		"success": {
			expRs: []model.Product{
				{
					ID:    1,
					Name:  "test1",
					Price: 1,
				},
				{
					ID:    2,
					Name:  "test2",
					Price: 2,
				},
			},
		},
		"empty": {
			isEmpty: true,
			expRs:   []model.Product{},
		},
		"error: db failed": {
			expDBFailed: true,
			expErr:      errors.New("models: failed to assign all query results to Product slice: bind failed to execute query: sql: database is closed"),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			TestWithTxDB(t, func(tx db.ContextExecutor) {
				// Given
				repo := New(tx)
				if tc.expDBFailed {
					dbMock, _, _ := sqlmock.New()
					dbMock.Close()
					repo = New(dbMock)
				}
				if !tc.isEmpty {
					testdata.LoadTestSQLFile(t, tx, "testdata/get_all_products.sql")
				}

				// When
				result, err := repo.GetAll(ctx)

				// Then
				if tc.expErr != nil {
					require.EqualError(t, err, tc.expErr.Error())
				} else {
					require.NoError(t, err)
					if !cmp.Equal(tc.expRs, result,
						cmpopts.IgnoreFields(model.Product{}, "CreatedAt", "UpdatedAt", "DeletedAt")) {
						t.Errorf("\n order mismatched. \n expected: %+v \n got: %+v \n diff: %+v", tc.expRs, result,
							cmp.Diff(tc.expRs, result, cmpopts.IgnoreFields(model.Product{}, "CreatedAt", "UpdatedAt", "DeletedAt")))
						t.FailNow()
					}
				}
			})
		})
	}
}

func TestImpl_DeleteById(t *testing.T) {
	type args struct {
		givenID     int64
		expDBFailed bool
		expErr      error
	}

	tcs := map[string]args{
		"success": {
			givenID: 1,
			expErr:  nil,
		},
		"error: not found": {
			givenID: 1000,
			expErr:  sql.ErrNoRows,
		},
		"error: db failed": {
			givenID:     1,
			expDBFailed: true,
			expErr:      errors.New("models: unable to select from product: bind failed to execute query: sql: database is closed"),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			TestWithTxDB(t, func(tx db.ContextExecutor) {
				// Given
				repo := New(tx)
				if tc.expDBFailed {
					dbMock, _, _ := sqlmock.New()
					dbMock.Close()
					repo = New(dbMock)
				}
				testdata.LoadTestSQLFile(t, tx, "testdata/delete_product_by_id.sql")

				// When
				err := repo.Delete(ctx, tc.givenID)

				// Then
				if tc.expErr != nil {
					require.EqualError(t, err, tc.expErr.Error())
				} else {
					require.NoError(t, err)
				}
			})
		})
	}
}
