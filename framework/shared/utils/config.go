package utils

import frameworkruntime "github.com/khiemnd777/noah_framework/runtime"

func LoadConfig[T any](path string) (*T, error) {
	return frameworkruntime.LoadYAML[T](path)
}
