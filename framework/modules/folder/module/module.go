package module

import (
	"database/sql"
	"time"

	frameworkhandler "github.com/khiemnd777/noah_framework/modules/folder/handler"
	frameworkrepository "github.com/khiemnd777/noah_framework/modules/folder/repository"
	frameworkservice "github.com/khiemnd777/noah_framework/modules/folder/service"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

type Config struct {
	Cache CacheConfig
}

type CacheConfig struct {
	ShortTTL time.Duration
	LongTTL  time.Duration
}

type Options struct {
	Database *sql.DB
	Config   Config
}

func Register(router frameworkhttp.Router, opts Options) error {
	repo := frameworkrepository.New(opts.Database)
	svc := frameworkservice.New(repo, frameworkservice.CacheConfig{
		ShortTTL: opts.Config.Cache.ShortTTL,
		LongTTL:  opts.Config.Cache.LongTTL,
	})
	frameworkhandler.New(svc).RegisterRoutes(router)
	return nil
}
