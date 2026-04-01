package app

import (
	"fmt"

	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

func ParseBody[T any](c frameworkhttp.Context) (*T, error) {
	var body T
	if err := c.BodyParser(&body); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}
	return &body, nil
}
