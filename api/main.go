package main

import (
	"fmt"
	"log"
	"os"

	entsql "entgo.io/ent/dialect/sql"
	gateway "github.com/khiemnd777/noah_api/gateway/runtime"
	"github.com/khiemnd777/noah_api/shared/bootstrap"
	"github.com/khiemnd777/noah_api/shared/cache"
	"github.com/khiemnd777/noah_api/shared/circuitbreaker"
	"github.com/khiemnd777/noah_api/shared/config"
	"github.com/khiemnd777/noah_api/shared/cron"
	"github.com/khiemnd777/noah_api/shared/db"
	"github.com/khiemnd777/noah_api/shared/db/ent"
	"github.com/khiemnd777/noah_api/shared/db/ent/generated"
	"github.com/khiemnd777/noah_api/shared/gen"
	"github.com/khiemnd777/noah_api/shared/logger"
	"github.com/khiemnd777/noah_api/shared/redis"
	"github.com/khiemnd777/noah_api/shared/utils"
	"github.com/khiemnd777/noah_api/shared/worker"
	_ "github.com/khiemnd777/noah_api/starter"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkcache "github.com/khiemnd777/noah_framework/pkg/cache"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

func main() {
	if err := config.EnsureEnvLoaded(); err != nil {
		log.Println("❌ Load .env failed:", err)
		os.Exit(1)
	}

	logger.Init()
	logger.SetComponent("api")

	log.Println("🔧 Loading config file...")
	if err := config.Init(utils.GetAppConfigPath()); err != nil {
		log.Println("❌ Load config failed:", err)
		os.Exit(1)
	}

	log.Println("✅ Config file loaded!")
	logger.Configure(logger.Options{
		ServiceName:  config.Get().Observability.ServiceName,
		Environment:  config.Get().Observability.Environment,
		Level:        config.Get().Observability.Logs.Level,
		RedactFields: config.Get().Observability.Logs.RedactFields,
		Component:    "api",
	})

	log.Println("🚀 Starting Project...",
		"project:", config.Get().Project.Name,
		"api_prefix:",
		config.Get().Project.BaseAPIPrefix,
		config.Get().Project.Version,
	)

	dbCfg := config.Get().Database

	dbClient, err := db.NewDatabaseClient(dbCfg)
	if err != nil {
		log.Fatalf("Cannot initialize DB: %v", err)
	}
	defer dbClient.Close()

	if err := dbClient.Connect(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	log.Println("Connected to DB successfully!")

	cache.InitTTLConstants()

	if err := frameworkruntime.ConfigureDefaultCache(toFrameworkCacheConfig(config.Get().Redis)); err != nil {
		log.Fatalf("Cannot initialize framework cache: %v", err)
	}

	// Initialize Redis
	redis.Init()

	if err := gen.GenerateEntClient(); err != nil {
		os.Exit(1)
	}

	sqlDB := dbClient.GetSQL() // Returns *sql.DB if Postgres, but nil Mongo
	_, entErr := ent.EntBootstrap(dbCfg.Provider, sqlDB, func(drv *entsql.Driver) any {
		return generated.NewClient(generated.Driver(drv))
	}, dbCfg.AutoMigrate)
	if entErr != nil {
		log.Fatalf("❌ Failed to init Ent client: %v", entErr)
		os.Exit(1)
	}

	if err := bootstrap.ApplySQLMigrations(sqlDB); err != nil {
		os.Exit(1)
	}

	if err := bootstrap.EnsureBaseRolesAndPermissions(dbCfg); err != nil {
		log.Fatalf("❌ Failed to seed base roles: %v", err)
	}

	circuitbreaker.Init()

	defer worker.StopAllWorkers()

	cron.StartAllCrons()

	app := frameworkruntime.NewApplication(frameworkapp.Config{
		BodyLimitMB: config.Get().Server.BodyLimitMB,
		Host:        config.Get().Server.Host,
		Port:        config.Get().Server.Port,
	})
	if err := gateway.Start(app); err != nil {
		log.Fatalf("Gateway error: %v", err)
	}
	if err := app.Listen(fmt.Sprintf("%s:%d", config.Get().Server.Host, config.Get().Server.Port)); err != nil {
		log.Fatalf("Gateway listen error: %v", err)
	}
}

func toFrameworkCacheConfig(cfg config.RedisConfig) frameworkcache.Config {
	instances := make(map[string]frameworkcache.InstanceConfig, len(cfg.Instances))
	for name, instance := range cfg.Instances {
		instances[name] = frameworkcache.InstanceConfig{
			DB:       instance.DB,
			Host:     instance.Host,
			Password: instance.Password,
			Port:     instance.Port,
		}
	}
	return frameworkcache.Config{
		DefaultInstance: "cache",
		Instances:       instances,
	}
}
