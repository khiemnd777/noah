package db

import (
	"database/sql"

	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

func SQLDB(client frameworkdb.Client) (*sql.DB, error) {
	return frameworkruntime.SQLDB(client)
}

func MustSQLDB(client frameworkdb.Client) *sql.DB {
	return frameworkruntime.MustSQLDB(client)
}
