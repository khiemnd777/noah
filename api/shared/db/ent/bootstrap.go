package ent

import (
	"context"
	"database/sql"
	"fmt"

	shareddb "github.com/khiemnd777/noah_api/shared/db"
	"github.com/khiemnd777/noah_api/shared/db/ent/generated"
	"github.com/khiemnd777/noah_api/shared/logger"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
)

func InitEntClientFromDatabase(client frameworkdb.Client, builder EntClientBuilder) (any, error) {
	sqlDB, err := shareddb.SQLDB(client)
	if err != nil {
		return nil, err
	}
	return InitEntClient(client.Provider(), sqlDB, builder)
}

func EntBootstrapFromDatabase(client frameworkdb.Client, builder EntClientBuilder, autoMigrate bool) (any, error) {
	sqlDB, err := shareddb.SQLDB(client)
	if err != nil {
		return nil, err
	}
	return EntBootstrap(client.Provider(), sqlDB, builder, autoMigrate)
}

func EntBootstrap(provider string, db *sql.DB, builder EntClientBuilder, autoMigrate bool) (any, error) {
	rawClient, err := InitEntClient(provider, db, builder)
	if err != nil {
		return nil, fmt.Errorf("init ent client failed: %w", err)
	}

	client := rawClient.(*generated.Client)
	// defer client.Close()

	if autoMigrate {
		if err := client.Schema.Create(context.Background()); err != nil {
			return nil, fmt.Errorf("auto create schema failed: %w", err)
		}
		logger.Info("📦 Ent schema created successfully")
	}

	return rawClient, nil
}
