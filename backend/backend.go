package backend

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/welaw/welaw/backend/database"
)

func NewTestDatabase(logger log.Logger) database.Database {
	connStr := os.Getenv("POSTGRES_CONNECTION_TEST")
	return database.NewDatabase(connStr, logger, &database.DatabaseConfigOptions{})
}

func NewDatabase(connStr string, logger log.Logger, opts *database.DatabaseConfigOptions) database.Database {
	return database.NewDatabase(connStr, logger, opts)
}
