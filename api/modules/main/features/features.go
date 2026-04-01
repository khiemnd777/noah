package features

import (
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

func Register(router frameworkhttp.Router) {
	router.Get("/", func(c frameworkhttp.Context) error {
		return c.JSON(map[string]any{
			"name":     "noah api sample",
			"status":   "ok",
			"features": []string{"sample-health", "sample-info"},
		})
	})

	router.Get("/health", func(c frameworkhttp.Context) error {
		return c.JSON(map[string]string{
			"status": "ok",
		})
	})
}
