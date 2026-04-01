package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/khiemnd777/noah_api/modules/main/features"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

const (
	defaultHost = "127.0.0.1"
	defaultPort = 8080
)

func main() {
	host := envOrDefault("API_HOST", defaultHost)
	port := envInt("API_PORT", defaultPort)

	app := frameworkruntime.NewApplication(frameworkapp.Config{
		Host: host,
		Port: port,
	})
	features.Register(app.Router())

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("sample api listening on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("listen failed: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
