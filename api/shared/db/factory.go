package db

import (
	"fmt"

	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"

	"github.com/khiemnd777/noah_api/shared/config"
)

func NewDatabaseClient(cfg config.DatabaseConfig) (frameworkdb.Client, error) {
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
	return client, nil
}
