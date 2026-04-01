package runtime

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	frameworklifecycle "github.com/khiemnd777/noah_framework/pkg/lifecycle"
	frameworkmodule "github.com/khiemnd777/noah_framework/pkg/module"
	"gopkg.in/yaml.v3"
)

type RunningModule = frameworklifecycle.ModuleInfo
type Registry = frameworklifecycle.Registry

func LoadRegistry(path string) (Registry, error) {
	store := ConfigureDefaultLifecycleStore(path)
	return store.Load()
}

func SaveRegistry(path string, reg Registry) error {
	store := ConfigureDefaultLifecycleStore(path)
	return store.Save(reg)
}

func UpdateRegistry(path string, update func(Registry)) error {
	store := ConfigureDefaultLifecycleStore(path)
	return store.Update(update)
}

func RegisterModule(path, name, route, host string, port int, external bool) error {
	store := ConfigureDefaultLifecycleStore(path)
	frameworklifecycle.SetDefaultStore(store)
	return frameworklifecycle.Register(name, route, host, port, external)
}

func DefaultDiscoveryRoots() []frameworkmodule.DiscoveryRoot {
	apiRoot := APIPath()

	return []frameworkmodule.DiscoveryRoot{
		{Name: "framework", Path: filepath.Join(FindRepoRoot(), "framework", "modules")},
		{Name: "api-main", Path: filepath.Join(apiRoot, "modules", "main")},
		{Name: "api", Path: filepath.Join(apiRoot, "modules")},
	}
}

func DiscoverModuleDescriptors(roots []frameworkmodule.DiscoveryRoot) ([]frameworkmodule.Descriptor, error) {
	if len(roots) == 0 {
		roots = DefaultDiscoveryRoots()
	}
	return DiscoverModules(roots)
}

func ResolveModuleDescriptor(moduleName string, roots []frameworkmodule.DiscoveryRoot) (frameworkmodule.Descriptor, error) {
	descriptors, err := DiscoverModuleDescriptors(roots)
	if err != nil {
		return frameworkmodule.Descriptor{}, err
	}

	for _, descriptor := range descriptors {
		if descriptor.ID == moduleName || descriptor.Name == moduleName {
			return descriptor, nil
		}
	}

	return frameworkmodule.Descriptor{}, fmt.Errorf("module %q not found in discovery roots", moduleName)
}

func ModuleConfigPath(moduleName string, roots []frameworkmodule.DiscoveryRoot) (string, error) {
	descriptor, err := ResolveModuleDescriptor(moduleName, roots)
	if err != nil {
		return "", err
	}
	return descriptor.ConfigPath, nil
}

func ModuleEntrypointPath(moduleName string, roots []frameworkmodule.DiscoveryRoot) (string, error) {
	apiEntryPath := APIModulePath(moduleName, "main.go")
	if _, err := os.Stat(apiEntryPath); err == nil {
		return apiEntryPath, nil
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	descriptor, err := ResolveModuleDescriptor(moduleName, roots)
	if err != nil {
		return "", err
	}
	return descriptor.EntryPath, nil
}

func GenerateRegistry(path string, roots []frameworkmodule.DiscoveryRoot, basePort int) (Registry, []*Reserved, error) {
	reg := Registry{}
	var reserved []*Reserved

	descriptors, err := DiscoverModuleDescriptors(roots)
	if err != nil {
		return nil, nil, err
	}

	for _, descriptor := range descriptors {
		_, route, host, port, external, err := loadServerSection(descriptor.ConfigPath)
		if err != nil {
			continue
		}

		destPort := basePort + port
		r, err := EnsurePortAvailable(host, destPort)
		if err != nil {
			continue
		}

		reserved = append(reserved, r)
		reg[descriptor.ID] = RunningModule{
			PID:      0,
			Host:     host,
			Port:     r.Port,
			Route:    route,
			RunAt:    time.Now(),
			External: external,
		}
	}

	if err := SaveRegistry(path, reg); err != nil {
		return nil, nil, err
	}

	return reg, reserved, nil
}

func loadServerSection(cfgPath string) (
	rawFile []byte, route, host string, port int, external bool, err error,
) {
	data, err := ReadExpandedYAML(cfgPath)
	if err != nil {
		return nil, "", "", 0, false, err
	}

	var raw struct {
		Server struct {
			Host  string `yaml:"host"`
			Port  int    `yaml:"port"`
			Route string `yaml:"route"`
		} `yaml:"server"`
		External bool `yaml:"external"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, "", "", 0, false, err
	}

	return data, raw.Server.Route, raw.Server.Host, raw.Server.Port, raw.External, nil
}
