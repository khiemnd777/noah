package runtime

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkcache "github.com/khiemnd777/noah_framework/pkg/cache"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
	frameworklifecycle "github.com/khiemnd777/noah_framework/pkg/lifecycle"
)

type ModuleDeps[T any] struct {
	AppConfig *AppConfig
	Config    *T
	DBClient  frameworkdb.Client
	DB        *sql.DB
	App       frameworkapp.Application
}

type ModuleOptions[T any] struct {
	Name        string
	ConfigPath  string
	GetServer   func(cfg *T) ServerConfig
	GetDatabase func(cfg *T) DatabaseConfig
	OnRegister  func(app frameworkapp.Application, deps *ModuleDeps[T]) error
}

func StartModule[T any](opts ModuleOptions[T]) error {
	appCfg, err := LoadYAML[AppConfig](APIPath("config.yaml"))
	if err != nil {
		return fmt.Errorf("load app config: %w", err)
	}

	moduleCfg, err := LoadYAML[T](opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("load module config: %w", err)
	}

	if err := ConfigureDefaultCache(toFrameworkCacheConfig(appCfg.Redis)); err != nil {
		return fmt.Errorf("configure cache: %w", err)
	}

	dbClient, err := NewDatabaseClient(toFrameworkDBConfig(opts.GetDatabase(moduleCfg)))
	if err != nil {
		return fmt.Errorf("create db client: %w", err)
	}
	defer dbClient.Close()

	if err := dbClient.Connect(); err != nil {
		return fmt.Errorf("connect db: %w", err)
	}

	appRuntime := NewApplication(frameworkapp.Config{
		BodyLimitMB: appCfg.Server.BodyLimitMB,
		Host:        opts.GetServer(moduleCfg).Host,
		Port:        opts.GetServer(moduleCfg).Port,
	})
	appRuntime.Router().Get("/health", func(c frameworkhttp.Context) error {
		return c.SendStatus(200)
	})

	deps := &ModuleDeps[T]{
		AppConfig: appCfg,
		Config:    moduleCfg,
		DBClient:  dbClient,
		DB:        MustSQLDB(dbClient),
		App:       appRuntime,
	}

	if err := opts.OnRegister(appRuntime, deps); err != nil {
		return err
	}

	return serveModuleApp(opts.Name, appRuntime)
}

func serveModuleApp(moduleName string, appRuntime frameworkapp.Application) error {
	storePath := APIPath("tmp", "runtime.json")
	ConfigureDefaultLifecycleStore(storePath)

	store, err := frameworklifecycle.DefaultStore()
	if err != nil {
		return fmt.Errorf("open lifecycle store: %w", err)
	}

	reg, err := store.Load()
	if err != nil {
		return fmt.Errorf("load runtime registry: %w", err)
	}

	rm, ok := reg[moduleName]
	if !ok || rm.Host == "" || rm.Port == 0 {
		return fmt.Errorf("runtime entry for %q not found", moduleName)
	}

	addr := fmt.Sprintf("%s:%d", rm.Host, rm.Port)
	listener, err := reservedListener(rm.Host, rm.Port)
	if err != nil {
		return err
	}

	if err := store.Update(func(reg frameworklifecycle.Registry) {
		current := reg[moduleName]
		current.PID = os.Getpid()
		current.RunAt = time.Now()
		current.Host = rm.Host
		current.Port = rm.Port
		reg[moduleName] = current
	}); err != nil {
		log.Printf("module %s: update runtime registry: %v", moduleName, err)
	}

	log.Printf("module %s listening on %s", moduleName, addr)
	return appRuntime.Serve(listener)
}

func toFrameworkDBConfig(cfg DatabaseConfig) frameworkdb.Config {
	return frameworkdb.Config{
		Provider: cfg.Provider,
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
	}
}

func toFrameworkCacheConfig(cfg RedisConfig) frameworkcache.Config {
	instances := make(map[string]frameworkcache.InstanceConfig, len(cfg.Instances))
	for name, instance := range cfg.Instances {
		instances[name] = frameworkcache.InstanceConfig{
			DB:        instance.DB,
			Host:      instance.Host,
			Password:  instance.Password,
			Port:      instance.Port,
			Username:  instance.Username,
			IsCluster: instance.IsCluster,
			UseTLS:    instance.UseTLS,
		}
	}

	return frameworkcache.Config{
		DefaultInstance: "cache",
		Instances:       instances,
	}
}
