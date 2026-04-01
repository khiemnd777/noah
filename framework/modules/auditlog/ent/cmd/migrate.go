package main

import (
	"fmt"
	"log"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/khiemnd777/noah_framework/modules/auditlog/config"
	"github.com/khiemnd777/noah_framework/modules/auditlog/ent/bootstrap"
	"github.com/khiemnd777/noah_framework/modules/auditlog/ent/generated"
	sharedConfig "github.com/khiemnd777/noah_framework/internal/legacy/shared/config"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/db"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/logger"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/utils"
	_ "github.com/lib/pq"
)

func main() {
	sharedConfig.Init(utils.GetAppConfigPath())
	cfg, _ := utils.LoadConfig[config.ModuleConfig](utils.GetModuleConfigPath("auditlog"))
	dbCfg := any(cfg).(interface {
		GetDatabase() sharedConfig.DatabaseConfig
	}).GetDatabase()

	dbClient, err := db.NewDatabaseClient(dbCfg)
	if err != nil {
		logger.Error(fmt.Sprintf("❌ Cannot create database client: %v", err))
	}
	defer dbClient.Close()

	if err := dbClient.Connect(); err != nil {
		logger.Error(fmt.Sprintf("❌ Failed to connect to database: %v", err))
	}

	sqlDB := db.MustSQLDB(dbClient)
	_, err = bootstrap.EntBootstrap(dbCfg.Provider, sqlDB, func(drv *entsql.Driver) any {
		return generated.NewClient(generated.Driver(drv))
	}, cfg.Database.AutoMigrate)
	if err != nil {
		logger.Error(fmt.Sprintf("❌ Failed to init Ent client: %v", err))
	}

	log.Println("✅ Migration completed successfully")
}
