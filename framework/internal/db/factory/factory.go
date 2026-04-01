package factory

import (
	"fmt"

	"github.com/khiemnd777/noah_framework/internal/db/driver"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
)

func NewClient(cfg frameworkdb.Config) (frameworkdb.Client, error) {
	switch cfg.Provider {
	case "postgres":
		return driver.NewPostgresClient(cfg.Postgres), nil
	case "mongodb":
		return driver.NewMongoClient(cfg.MongoDB), nil
	default:
		return nil, fmt.Errorf("unsupported database provider: %s", cfg.Provider)
	}
}
