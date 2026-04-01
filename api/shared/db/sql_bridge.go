package db

import (
	"database/sql"
	"fmt"

	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
)

func SQLDB(client frameworkdb.Client) (*sql.DB, error) {
	bridge, ok := client.(frameworkdb.SQLBridge)
	if !ok {
		return nil, fmt.Errorf("database provider %q does not expose an sql bridge", client.Provider())
	}

	sqlDB, ok := bridge.SQLDB().(*sql.DB)
	if !ok || sqlDB == nil {
		return nil, fmt.Errorf("database provider %q returned an invalid sql bridge", client.Provider())
	}

	return sqlDB, nil
}

func MustSQLDB(client frameworkdb.Client) *sql.DB {
	sqlDB, err := SQLDB(client)
	if err != nil {
		panic(err)
	}
	return sqlDB
}
