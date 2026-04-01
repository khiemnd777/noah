package utils

import (
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
	"os"
)

func DiscoverAllModules() ([]string, error) {
	descriptors, err := frameworkruntime.DiscoverModuleDescriptors(nil)
	if err != nil {
		return nil, err
	}
	modules := make([]string, 0, len(descriptors))
	for _, descriptor := range descriptors {
		modules = append(modules, descriptor.ID)
	}
	return modules, nil
}

func GetModuleRoute(defaultRoute string) string {
	if os.Getenv("GATEWAY_MODE") == "true" {
		return "/"
	}
	return defaultRoute
}
