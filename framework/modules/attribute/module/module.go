package module

import (
	"context"

	frameworkhandler "github.com/khiemnd777/noah_framework/modules/attribute/handler"
	frameworkrepository "github.com/khiemnd777/noah_framework/modules/attribute/repository"
	frameworkservice "github.com/khiemnd777/noah_framework/modules/attribute/service"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

type Config struct {
	AutoMigrate bool
	Cache       frameworkservice.CacheConfig
}

type Options struct {
	Config   Config
	Database frameworkdb.Client
}

func Register(router frameworkhttp.Router, opts Options) error {
	repo, err := frameworkrepository.New(opts.Database)
	if err != nil {
		return err
	}

	if opts.Config.AutoMigrate {
		if err := repo.EnsureSchema(context.Background()); err != nil {
			return err
		}
	}

	attributeService := frameworkservice.New(repo, opts.Config.Cache)
	attributeHandler := frameworkhandler.NewAttributeHandler(attributeService)
	attributeHandler.RegisterRoutes(router)

	optionService := frameworkservice.NewOptionService(repo, opts.Config.Cache)
	optionHandler := frameworkhandler.NewAttributeOptionHandler(optionService)
	optionHandler.RegisterRoutes(router)
	return nil
}
