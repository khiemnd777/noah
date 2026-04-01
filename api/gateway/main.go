package main

import (
	"fmt"
	"log"
	"os"

	gateway "github.com/khiemnd777/noah_api/gateway/runtime"
	"github.com/khiemnd777/noah_api/shared/config"
	"github.com/khiemnd777/noah_api/shared/logger"
	"github.com/khiemnd777/noah_api/shared/utils"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

func main() {
	if err := config.EnsureEnvLoaded(); err != nil {
		log.Println("❌ Load .env failed:", err)
		os.Exit(1)
	}

	logger.Init()
	logger.SetComponent("gateway")

	if err := config.Init(utils.GetAppConfigPath()); err != nil {
		log.Println("❌ Load config failed:", err)
		os.Exit(1)
	}
	logger.Configure(logger.Options{
		ServiceName:  config.Get().Observability.ServiceName,
		Environment:  config.Get().Observability.Environment,
		Level:        config.Get().Observability.Logs.Level,
		RedactFields: config.Get().Observability.Logs.RedactFields,
		Component:    "gateway",
	})

	logger.Info("Starting API Gateway...")
	app := frameworkruntime.NewApplication(frameworkapp.Config{
		BodyLimitMB: config.Get().Server.BodyLimitMB,
		Host:        config.Get().Server.Host,
		Port:        config.Get().Server.Port,
	})

	if err := gateway.Start(app); err != nil {
		logger.Error("Gateway error", err)
	}

	if err := app.Listen(fmt.Sprintf("%s:%d", config.Get().Server.Host, config.Get().Server.Port)); err != nil {
		logger.Error("Gateway listen error", err)
	}
}
