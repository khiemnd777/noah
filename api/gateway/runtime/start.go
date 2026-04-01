package gateway

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/khiemnd777/noah_api/gateway/proxy"
	"github.com/khiemnd777/noah_api/scripts/module_runner/runner"
	"github.com/khiemnd777/noah_api/shared/config"
	"github.com/khiemnd777/noah_api/shared/logger"
	"github.com/khiemnd777/noah_api/shared/runtime"
	"github.com/khiemnd777/noah_api/shared/utils"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

func Start(app frameworkapp.Application) error {
	reg, reserved, err := runtime.GenerateRegistry(utils.GetFullPath("modules"))
	if err != nil {
		return fmt.Errorf("failed to load modules: %w", err)
	}

	for _, r := range reserved {
		_ = r.Listener.Close()
	}

	var wg sync.WaitGroup

	for name, m := range reg {
		wg.Add(1)

		go func(name string, m runtime.RunningModule) {
			defer wg.Done()

			logger.Info(fmt.Sprintf("Launching module [%s]...", name))
			// boot.TryStartModule(m.Name, m.Port)
			// runner.StartModulesInBatch([]string{m.Name})
			if err := runner.StartModule(name); err != nil {
				logger.Warn(fmt.Sprintf("Module [%s] failed to start", name), err)
				return
			}

			logger.Info(fmt.Sprintf("Module [%s] is ready", name))

			if m.External {
				target := fmt.Sprintf("http://%s:%d", m.Host, m.Port)
				logger.Info(fmt.Sprintf("Registering route %s → %s", m.Route, target))
				err := proxy.RegisterReverseProxy(app, m.Route, []string{target})
				if err != nil {
					logger.Warn(fmt.Sprintf("Failed to register module [%s]", name), err)
				}
			} else {
				logger.Info(fmt.Sprintf("Module [%s] is internal, skipping route registration", name))
			}
		}(name, m)
	}

	wg.Wait()

	runner.SyncRunningModules()

	srvCfg := config.Get().Server
	addr := fmt.Sprintf("%s:%d", srvCfg.Host, srvCfg.Port)
	logger.Info("Gateway listening on " + addr)

	// Get dev.log
	routeGetLog(app.Router())

	// API Home Page
	app.Router().Get("/", func(c frameworkhttp.Context) error {
		return c.SendString(`🚀 NOAH API has launched already!`)
	})
	return nil
}

func routeGetLog(router frameworkhttp.Router) {
	router.Get("/__log", func(c frameworkhttp.Context) error {
		if c.Get("Authorization") != fmt.Sprintf("Bearer %s", config.Get().Auth.InternalLogToken) {
			return c.Status(403).SendString("Access denied")
		}

		file, err := os.Open("dev.log")

		if err != nil {
			return c.Status(500).SendString("Failed to open log")
		}
		defer file.Close()

		lines := []string{}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if len(lines) > 100 {
			lines = lines[len(lines)-100:]
		}

		return c.SendString(strings.Join(lines, "\n"))
	})
}
