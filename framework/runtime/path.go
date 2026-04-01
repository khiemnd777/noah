package runtime

import (
	"os"
	"path/filepath"
)

func FindRepoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			return dir
		}
		if _, err := os.Stat(filepath.Join(dir, "framework", "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return dir
		}
		dir = parent
	}
}

func APIPath(parts ...string) string {
	root := FindRepoRoot()
	all := append([]string{root, "api"}, parts...)
	return filepath.Join(all...)
}

func FrameworkModulePath(moduleName string, parts ...string) string {
	root := FindRepoRoot()
	all := append([]string{root, "framework", "modules", moduleName}, parts...)
	return filepath.Join(all...)
}
