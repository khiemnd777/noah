package runtime

import (
	"database/sql"
	"fmt"

	internalfactory "github.com/khiemnd777/noah_framework/internal/db/factory"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
)

func NewDatabaseClient(cfg frameworkdb.Config) (frameworkdb.Client, error) {
	return internalfactory.NewClient(cfg)
}

type sqlBridge interface {
	SQLDB() *sql.DB
}

func SQLDB(client frameworkdb.Client) (*sql.DB, error) {
	bridge, ok := client.(sqlBridge)
	if !ok {
		return nil, fmt.Errorf("database provider %q does not expose an sql bridge", client.Provider())
	}

	sqlDB := bridge.SQLDB()
	if sqlDB == nil {
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
