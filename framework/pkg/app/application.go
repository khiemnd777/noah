package app

import frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"

type Config struct {
	Host        string
	Port        int
	BodyLimitMB int
}

type Application interface {
	Native() any
	Router() frameworkhttp.Router
	Listen(addr string) error
	Serve(listener any) error
	Run() error
}
