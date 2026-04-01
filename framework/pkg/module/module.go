package module

import frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"

type Definition interface {
	Name() string
	Register(app frameworkapp.Application) error
}

// DiscoveryRoot defines a filesystem root that may contain runnable modules.
// Roots are ordered; earlier roots win when the same module path is discovered
// through overlapping roots.
type DiscoveryRoot struct {
	Name string
	Path string
}

// Descriptor describes a discovered runnable module directory.
type Descriptor struct {
	ID         string
	Name       string
	RootName   string
	RootPath   string
	Path       string
	ConfigPath string
	EntryPath  string
}
