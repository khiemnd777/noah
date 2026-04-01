package utils

import (
	"strconv"

	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

func GetParamAsNillableInt(c frameworkhttp.Context, paramName string) (*int, error) {
	return frameworkruntime.ParamAsNullableInt(c, paramName)
}

func GetParamAsInt(c frameworkhttp.Context, paramName string) (int, error) {
	return strconv.Atoi(c.Params(paramName))
}

func GetParamAsString(c frameworkhttp.Context, paramName string) string {
	return c.Params(paramName)
}

func GetQueryAsNillableInt(c frameworkhttp.Context, queryName string, defaultValue ...string) (*int, error) {
	return frameworkruntime.QueryAsNullableInt(c, queryName, defaultValue...)
}

func GetQueryAsInt(c frameworkhttp.Context, queryName string, defaultValue ...int) int {
	return frameworkruntime.QueryAsInt(c, queryName, defaultValue...)
}

func GetQueryAsInt64(c frameworkhttp.Context, queryName string, defaultValue ...int64) int64 {
	return frameworkruntime.QueryAsInt64(c, queryName, defaultValue...)
}

func GetQueryAsString(c frameworkhttp.Context, queryName string, defaultValue ...string) string {
	return c.Query(queryName, defaultValue...)
}

func GetQueryAsFloat64Pointer(c frameworkhttp.Context, name string) (*float64, error) {
	return frameworkruntime.QueryAsFloat64Pointer(c, name)
}
