package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestImpl_GetByStatusAndLimit(t *testing.T) {
	type args struct {
		givenStatus enum.PaymentStatus
		expDBFailed bool
		givenLimit  int
		expRs       []model.Payment
		expErr      error
	}

	tcs := map[string]args{
		"success": {
			givenStatus: enum.PaymentStatusPending,
			givenLimit:  100,
			expRs: []model.Payment{
				{
					ID:                     1000,
					BillingOrderID:         "999999999",
					TransactionReferenceID: "txn_ref_id",
					IamID:                  "auth0|5d4b5e6f7g8h9i0j1k2l3m4",
					Amount:                 500,
					Points:                 500,
					Currency:               enum.CurrencySGD,
					Type:                   enum.PaymentTypeUtilitiesAdhocPay,
					Status:                 enum.PaymentStatusPending,
					TransactionID:          0,
				},
				{
					ID:                     1001,
					BillingOrderID:         "888888888",
					TransactionReferenceID: "txn_ref_id_2",
					IamID:                  "auth0|fdsfdsfdsfsdfsdfsdfdfgfg",
					Amount:                 500,
					Points:                 500,
					Currency:               enum.CurrencySGD,
					Type:                   enum.PaymentTypeUtilitiesAdhocPay,
					Status:                 enum.PaymentStatusPending,
					TransactionID:          0,
				},
			},
		},
		"empty": {
			givenStatus: enum.PaymentStatusCancelled,
			givenLimit:  100,
			expRs:       []model.Payment{},
		},
		"error: db failed": {
			givenStatus: enum.PaymentStatusPending,
			givenLimit:  100,
			expDBFailed: true,
			expErr:      errors.New("sql: database is closed"),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			connectionStr := "postgres://postgres:postgres@localhost:5432/scc-pg?sslmode=disable"
			db, err := sql.Open("postgres", connectionStr)
			require.NoError(t, err)

			tx, err := db.BeginTx(ctx, &sql.TxOptions{})
			require.NoError(t, err)
			defer tx.Rollback()

			tx.Exec("INSERT...")
			repo := New(tx)

			// When
			result, err := repo.GetByStatusAndLimit(ctx, tc.givenStatus, tc.givenLimit)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				if !cmp.Equal(tc.expRs, result,
					cmpopts.IgnoreFields(model.Payment{}, "CreatedAt", "UpdatedAt")) {
					t.Errorf("\n order mismatched. \n expected: %+v \n got: %+v \n diff: %+v", tc.expRs, result,
						cmp.Diff(tc.expRs, result, cmpopts.IgnoreFields(model.Payment{}, "CreatedAt", "UpdatedAt")))
					t.FailNow()
				}
			}

			// testpg.WithTxDB(t, func(tx pg.BeginnerExecutor) {
			// 	// Given
			// 	instance := New(tx)
			// 	if tc.expDBFailed {
			// 		dbMock := sqlmock.New()
			// 		dbMock.Close()
			// 		instance = New(dbMock)
			// 	}
			// 	LoadTestSQLFile(t, tx, "testdata/get_all_products.sql")
			// 	// When
			// 	result, err := instance.GetByStatusAndLimit(ctx, tc.givenStatus, tc.givenLimit)

			// 	// Then
			// 	if tc.expErr != nil {
			// 		require.EqualError(t, err, tc.expErr.Error())
			// 	} else {
			// 		require.NoError(t, err)
			// 		if !cmp.Equal(tc.expRs, result,
			// 			cmpopts.IgnoreFields(model.Payment{}, "CreatedAt", "UpdatedAt")) {
			// 			t.Errorf("\n order mismatched. \n expected: %+v \n got: %+v \n diff: %+v", tc.expRs, result,
			// 				cmp.Diff(tc.expRs, result, cmpopts.IgnoreFields(model.Payment{}, "CreatedAt", "UpdatedAt")))
			// 			t.FailNow()
			// 		}
			// 	}
			// })
		})
	}
}
