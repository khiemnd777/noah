package utils

import (
	"os"
	"path/filepath"

	frameworkmodule "github.com/khiemnd777/noah_framework/pkg/module"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

func DiscoverAllModules() ([]string, error) {
	apiRoot := GetProjectRootDir()
	repoRoot := filepath.Clean(filepath.Join(apiRoot, ".."))

	descriptors, err := frameworkruntime.DiscoverModules([]frameworkmodule.DiscoveryRoot{
		{Name: "framework", Path: filepath.Join(repoRoot, "framework", "modules")},
		{Name: "api-main", Path: filepath.Join(apiRoot, "modules", "main")},
		{Name: "api", Path: filepath.Join(apiRoot, "modules")},
	})
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
