package db

import (
	"database/sql"
	"fmt"

	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"

	"github.com/khiemnd777/noah_api/shared/config"
	dbinterface "github.com/khiemnd777/noah_api/shared/db/interface"
)

func NewDatabaseClient(cfg config.DatabaseConfig) (dbinterface.DatabaseClient, error) {
	client, err := frameworkruntime.NewDatabaseClient(frameworkdb.Config{
		Provider:    cfg.Provider,
		AutoMigrate: cfg.AutoMigrate,
		Postgres: frameworkdb.PostgresConfig{
			Host:     cfg.Postgres.Host,
			Port:     cfg.Postgres.Port,
			User:     cfg.Postgres.User,
			Password: cfg.Postgres.Password,
			Name:     cfg.Postgres.Name,
			SSLMode:  cfg.Postgres.SSLMode,
		},
		MongoDB: frameworkdb.MongoConfig{
			URI:      cfg.MongoDB.URI,
			Database: cfg.MongoDB.Database,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("unsupported database provider: %s", cfg.Provider)
	}
	return &frameworkClientAdapter{client: client}, nil
}

type frameworkClientAdapter struct {
	client frameworkdb.Client
}

func (a *frameworkClientAdapter) Connect() error {
	return a.client.Connect()
}

func (a *frameworkClientAdapter) Close() error {
	return a.client.Close()
}

func (a *frameworkClientAdapter) GetSQL() *sql.DB {
	bridge, ok := a.client.(frameworkdb.SQLBridge)
	if !ok {
		return nil
	}
	sqlDB, _ := bridge.SQLDB().(*sql.DB)
	return sqlDB
}
