package runtime

import (
	"strconv"

	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

const (
	maxInt = int(^uint(0) >> 1)
	minInt = -maxInt - 1
)

func ParamAsNullableInt(c frameworkhttp.Context, paramName string) (*int, error) {
	if raw := c.Params(paramName); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return nil, err
		}
		return &parsed, nil
	}
	return nil, nil
}

func QueryAsNullableInt(c frameworkhttp.Context, queryName string, defaultValue ...string) (*int, error) {
	if raw := c.Query(queryName, defaultValue...); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return nil, err
		}
		return &parsed, nil
	}
	return nil, nil
}

func QueryAsInt(c frameworkhttp.Context, queryName string, defaultValue ...int) int {
	raw := c.Query(queryName)
	if raw == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	parsed, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || parsed > int64(maxInt) || parsed < int64(minInt) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	return int(parsed)
}

func QueryAsInt64(c frameworkhttp.Context, queryName string, defaultValue ...int64) int64 {
	raw := c.Query(queryName)
	if raw == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	parsed, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	return parsed
}

func QueryAsFloat64Pointer(c frameworkhttp.Context, name string) (*float64, error) {
	val := c.Query(name)
	if val == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
