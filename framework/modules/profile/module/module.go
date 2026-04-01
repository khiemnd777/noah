package module

import (
	"database/sql"
	"time"

	frameworkhandler "github.com/khiemnd777/noah_framework/modules/profile/handler"
	frameworkrepository "github.com/khiemnd777/noah_framework/modules/profile/repository"
	frameworkservice "github.com/khiemnd777/noah_framework/modules/profile/service"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

type Config struct {
	Cache CacheConfig
}

type CacheConfig struct {
	MediumTTL time.Duration
}

type Options struct {
	Database *sql.DB
	Config   Config
}

func Register(router frameworkhttp.Router, opts Options) error {
	repo := frameworkrepository.New(opts.Database)
	svc := frameworkservice.New(repo, frameworkservice.CacheConfig{
		MediumTTL: opts.Config.Cache.MediumTTL,
	})
	frameworkhandler.New(svc).RegisterRoutes(router)
	return nil
}
