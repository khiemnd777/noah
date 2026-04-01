package bootstrap

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/khiemnd777/noah_api/modules/auditlog/ent/generated"
	"github.com/khiemnd777/noah_api/shared/db/ent"
	"github.com/khiemnd777/noah_api/shared/logger"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
)

func EntBootstrapFromDatabase(client frameworkdb.Client, builder ent.EntClientBuilder, autoMigrate bool) (any, error) {
	rawClient, err := ent.InitEntClientFromDatabase(client, builder)
	if err != nil {
		return nil, fmt.Errorf("init ent client failed: %w", err)
	}

	typed := rawClient.(*generated.Client)
	if autoMigrate {
		if err := typed.Schema.Create(context.Background()); err != nil {
			return nil, fmt.Errorf("auto create schema failed: %w", err)
		}
		logger.Info("📦 Ent schema created successfully")
	}

	return rawClient, nil
}

func EntBootstrap(provider string, db *sql.DB, builder ent.EntClientBuilder, autoMigrate bool) (any, error) {
	rawClient, err := ent.InitEntClient(provider, db, builder)
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
