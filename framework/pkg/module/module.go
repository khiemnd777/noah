package module

import frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"

type Definition interface {
	Name() string
	Register(app frameworkapp.Application) error
}
