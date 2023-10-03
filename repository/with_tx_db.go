package repository

import (
	"chi-demo/db"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithTxDB(t *testing.T, testFunc func(tx db.ContextExecutor)) {
	ctx := context.Background()

	db, err := db.Init()
	require.NoError(t, err)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	defer tx.Rollback()

	testFunc(tx)
}
