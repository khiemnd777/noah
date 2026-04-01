package runtime

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/khiemnd777/noah_api/shared/app"
	"github.com/khiemnd777/noah_api/shared/config"
	"github.com/khiemnd777/noah_api/shared/utils"
	frameworklifecycle "github.com/khiemnd777/noah_framework/pkg/lifecycle"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
	"gopkg.in/yaml.v3"
)

type RunningModule = frameworklifecycle.ModuleInfo
type Registry = frameworklifecycle.Registry

var (
	registryPath = "tmp/runtime.json"
	storeOnce    sync.Once
)

func lifecycleStore() frameworklifecycle.Store {
	storeOnce.Do(func() {
		frameworkruntime.ConfigureDefaultLifecycleStore(registryPath)
	})

	store, _ := frameworklifecycle.DefaultStore()
	return store
}

func LoadRegistry() (Registry, error) {
	return lifecycleStore().Load()
}

func SaveRegistry(reg Registry) error {
	return lifecycleStore().Save(reg)
}

func Register(name, route, host string, port int, external bool) error {
	return frameworklifecycle.Register(name, route, host, port, external)
}

func getDestPort(port int) int {
	mCfg, _ := utils.LoadConfig[config.AppConfig](utils.GetAppConfigPath())
	mPort := mCfg.Server.Port
	return mPort + port
}

func GenerateRegistry(modDir string) (Registry, []*app.Reserved, error) {
	reg := Registry{}
	var rs []*app.Reserved

	entries, err := os.ReadDir(modDir)
	if err != nil {
		return nil, nil, err
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()

		cfgFile := utils.GetModuleConfigPath(name)
		_, routeFromCfg, hostFromCfg, portFromCfg, externalFromCfg, err := loadServerSection(cfgFile)
		if err != nil {
			fmt.Printf("⚠️  skip %s: %v\n", name, err)
			continue
		}

		host := hostFromCfg
		port := getDestPort(portFromCfg)
		r, portErr := app.EnsurePortAvailable(host, port)
		if portErr != nil {
			fmt.Printf("🛑  cannot alloc port for %s\n", name)
			continue
		}

		rs = append(rs, r)
		reg[name] = RunningModule{
			PID:      0,
			Host:     host,
			Port:     r.Port,
			Route:    routeFromCfg,
			RunAt:    time.Now(),
			External: externalFromCfg,
		}
	}

	if err := SaveRegistry(reg); err != nil {
		return nil, nil, err
	}
	return reg, rs, nil
}

func UpdateRegistry(update func(Registry)) error {
	return lifecycleStore().Update(update)
}

func loadServerSection(cfgPath string) (
	rawFile []byte, route, host string, port int, external bool, err error,
) {
	data, err := config.ReadExpandedYAML(cfgPath)
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
