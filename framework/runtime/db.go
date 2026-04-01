package runtime

import (
	internalfactory "github.com/khiemnd777/noah_framework/internal/db/factory"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
)

func NewDatabaseClient(cfg frameworkdb.Config) (frameworkdb.Client, error) {
	return internalfactory.NewClient(cfg)
}
