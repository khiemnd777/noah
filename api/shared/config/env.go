package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/subosito/gotenv"
	"gopkg.in/yaml.v3"
)

var loadDotEnvOnce sync.Once

func EnsureEnvLoaded() error {
	var loadErr error

	loadDotEnvOnce.Do(func() {
		loadErr = ensureEnvLoaded()
	})

	return loadErr
}

func ensureEnvLoaded() error {
	start, err := os.Getwd()
	if err != nil {
		return err
	}

	if repoRoot, ok := findRepoRoot(start); ok {
		if err := loadEnvFromDir(repoRoot); err != nil {
			return err
		}
	}

	if apiRoot, ok := findAncestorNamed(start, "api"); ok {
		if err := loadEnvFromDir(apiRoot); err != nil {
			return err
		}
	} else {
		fallback := filepath.Join(start, "api")
		if stat, err := os.Stat(fallback); err == nil && stat.IsDir() {
			if err := loadEnvFromDir(fallback); err != nil {
				return err
			}
		}
	}

	applySharedEnvOverrides()

	return nil
}

func envFileCandidates() []string {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("APP_ENV")), "production") {
		return []string{".env.prod", ".env"}
	}
	return []string{".env"}
}

func loadEnvFromDir(dir string) error {
	for _, name := range envFileCandidates() {
		candidate := filepath.Join(dir, name)
		if _, err := os.Stat(candidate); err == nil {
			if err := gotenv.Load(candidate); err != nil {
				return fmt.Errorf("load .env %s: %w", candidate, err)
			}
			return nil
		}
	}

	return nil
}

func applySharedEnvOverrides() {
	if frontendOrigin := strings.TrimSpace(os.Getenv("APP_FE_ORIGIN")); frontendOrigin != "" {
		os.Setenv("M_MAIN_CLIENT_BASE_URL", frontendOrigin)
	}
}

func findRepoRoot(start string) (string, bool) {
	for current := start; ; current = filepath.Dir(current) {
		if dirExists(filepath.Join(current, "api")) && dirExists(filepath.Join(current, "fe")) {
			return current, true
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", false
		}
	}
}

func findAncestorNamed(start string, name string) (string, bool) {
	for current := start; ; current = filepath.Dir(current) {
		if filepath.Base(current) == name {
			return current, true
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", false
		}
	}
}

func dirExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func ReadExpandedYAML(path string) ([]byte, error) {
	if err := EnsureEnvLoaded(); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return []byte(os.ExpandEnv(string(data))), nil
}

func UnmarshalYAMLFile(path string, out any) error {
	data, err := ReadExpandedYAML(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, out); err != nil {
		return err
	}

	return nil
}

func NewExpandedYAMLReader(path string) (*bytes.Reader, string, error) {
	data, err := ReadExpandedYAML(path)
	if err != nil {
		return nil, "", err
	}

	configType := strings.TrimPrefix(filepath.Ext(path), ".")
	if configType == "" {
		configType = "yaml"
	}

	return bytes.NewReader(data), configType, nil
}
