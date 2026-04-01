package runtime

import (
	"github.com/khiemnd777/noah_api/shared/config"
	"github.com/khiemnd777/noah_api/shared/utils"
	frameworklifecycle "github.com/khiemnd777/noah_framework/pkg/lifecycle"
	frameworkmodule "github.com/khiemnd777/noah_framework/pkg/module"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

type RunningModule = frameworklifecycle.ModuleInfo
type Registry = frameworklifecycle.Registry

const registryPath = "tmp/runtime.json"

func LoadRegistry() (Registry, error) {
	return frameworkruntime.LoadRegistry(registryPath)
}

func SaveRegistry(reg Registry) error {
	return frameworkruntime.SaveRegistry(registryPath, reg)
}

func Register(name, route, host string, port int, external bool) error {
	return frameworkruntime.RegisterModule(registryPath, name, route, host, port, external)
}

func getDestPort(port int) int {
	mCfg, _ := utils.LoadConfig[config.AppConfig](utils.GetAppConfigPath())
	mPort := mCfg.Server.Port
	return mPort + port
}

func DefaultDiscoveryRoots() []frameworkmodule.DiscoveryRoot {
	return frameworkruntime.DefaultDiscoveryRoots()
}

func DiscoverModuleDescriptors() ([]frameworkmodule.Descriptor, error) {
	return frameworkruntime.DiscoverModuleDescriptors(DefaultDiscoveryRoots())
}

func ResolveModuleDescriptor(moduleName string) (frameworkmodule.Descriptor, error) {
	return frameworkruntime.ResolveModuleDescriptor(moduleName, DefaultDiscoveryRoots())
}

func ModuleConfigPath(moduleName string) (string, error) {
	return frameworkruntime.ModuleConfigPath(moduleName, DefaultDiscoveryRoots())
}

func ModuleEntrypointPath(moduleName string) (string, error) {
	return frameworkruntime.ModuleEntrypointPath(moduleName, DefaultDiscoveryRoots())
}

func GenerateRegistry(roots []frameworkmodule.DiscoveryRoot) (Registry, []*frameworkruntime.Reserved, error) {
	return frameworkruntime.GenerateRegistry(registryPath, roots, getDestPort(0))
}

func UpdateRegistry(update func(Registry)) error {
	return frameworkruntime.UpdateRegistry(registryPath, update)
}
