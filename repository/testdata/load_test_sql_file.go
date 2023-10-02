package testdata

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// LoadTestSQLFile load test sql data from a file
func LoadTestSQLFile(t *testing.T, tx *sql.Tx, filename string) {
	body, err := os.ReadFile(filename)
	require.NoError(t, err)

	_, err = tx.Exec(string(body))
	require.NoError(t, err)
}
