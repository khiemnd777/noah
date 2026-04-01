package search

import (
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/logger"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/middleware/rbac"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/modules/search"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/modules/search/model"
)

func init() {
	logger.Debug("[GuardSearch] Register Material")
	search.RegisterGuard("material", func(ctx search.GuardCtx, rows []model.Row) []model.Row {
		perms := ctx.Perms

		if !rbac.HasAnyPerm(perms, "material.search") {
			return []model.Row{}
		}

		out := make([]model.Row, 0, len(rows))
		out = append(out, rows...)

		return out
	})
}
