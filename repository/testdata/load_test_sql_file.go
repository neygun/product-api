package testdata

import (
	"chi-demo/db"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// LoadTestSQLFile load test sql data from a file
func LoadTestSQLFile(t *testing.T, tx db.ContextExecutor, filename string) {
	body, err := os.ReadFile(filename)
	require.NoError(t, err)

	_, err = tx.Exec(string(body))
	require.NoError(t, err)
}
