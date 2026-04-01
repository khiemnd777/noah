package runtime

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/subosito/gotenv"
	"gopkg.in/yaml.v3"
)

type ProjectConfig struct {
	Name          string `yaml:"name"`
	BaseAPIPrefix string `yaml:"baseapiprefix"`
	Version       string `yaml:"version"`
}

type ServerConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Route       string `yaml:"route"`
	BodyLimitMB int    `yaml:"body_limit_mb"`
}

type AuthConfig struct {
	Secret            string `yaml:"secret"`
	InternalAuthToken string `yaml:"internalauthtoken"`
	InternalLogToken  string `yaml:"internallogtoken"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

type MongoConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
}

type DatabaseConfig struct {
	Provider    string         `yaml:"provider"`
	AutoMigrate bool           `yaml:"automigrate"`
	Postgres    PostgresConfig `yaml:"postgres"`
	MongoDB     MongoConfig    `yaml:"mongodb"`
}

type RedisInstanceConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	DB        int    `yaml:"db"`
	IsCluster bool   `yaml:"iscluster"`
	UseTLS    bool   `yaml:"usetls"`
}

type CacheTTL struct {
	Short  time.Duration `yaml:"short"`
	Medium time.Duration `yaml:"medium"`
	Long   time.Duration `yaml:"long"`
	Static time.Duration `yaml:"static"`
}

type CacheConfig struct {
	TTL CacheTTL `yaml:"ttl"`
}

type RedisConfig struct {
	Instances map[string]RedisInstanceConfig `yaml:"instances"`
}

type CircuitBreakerConfig struct {
	Interval            time.Duration `yaml:"interval"`
	Timeout             time.Duration `yaml:"timeout"`
	ConsecutiveFailures int           `yaml:"consecutivefailures"`
}

type RetryConfig struct {
	MaxAttempts int           `yaml:"maxattempts"`
	Delay       time.Duration `yaml:"delay"`
}

type ObservabilityLogsConfig struct {
	Enabled      bool     `yaml:"enabled"`
	Level        string   `yaml:"level"`
	RedactFields []string `yaml:"redact_fields"`
}

type ObservabilityLokiConfig struct {
	BaseURL        string        `yaml:"base_url"`
	TenantID       string        `yaml:"tenant_id"`
	BearerToken    string        `yaml:"bearer_token"`
	Timeout        time.Duration `yaml:"timeout"`
	StreamSelector string        `yaml:"stream_selector"`
	MaxQueryLimit  int           `yaml:"max_query_limit"`
}

type ObservabilityConfig struct {
	ServiceName string                  `yaml:"service_name"`
	Environment string                  `yaml:"environment"`
	Logs        ObservabilityLogsConfig `yaml:"logs"`
	Loki        ObservabilityLokiConfig `yaml:"loki"`
}

type AppConfig struct {
	Project        ProjectConfig        `yaml:"project"`
	Server         ServerConfig         `yaml:"server"`
	Auth           AuthConfig           `yaml:"auth"`
	Database       DatabaseConfig       `yaml:"database"`
	Cache          CacheConfig          `yaml:"cache"`
	Redis          RedisConfig          `yaml:"redis"`
	CircuitBreaker CircuitBreakerConfig `yaml:"circuitbreaker"`
	Retry          RetryConfig          `yaml:"retry"`
	Observability  ObservabilityConfig  `yaml:"observability"`
}

var loadDotEnvOnce sync.Once

func EnsureEnvLoaded() error {
	var loadErr error

	loadDotEnvOnce.Do(func() {
		start, err := os.Getwd()
		if err != nil {
			loadErr = err
			return
		}

		root := start
		candidates := []string{".env"}
		if strings.EqualFold(strings.TrimSpace(os.Getenv("APP_ENV")), "production") {
			candidates = []string{".env.prod", ".env"}
		}

		for {
			for _, name := range candidates {
				candidate := filepath.Join(root, name)
				if _, err := os.Stat(candidate); err == nil {
					if err := gotenv.Load(candidate); err != nil {
						loadErr = fmt.Errorf("load env %s: %w", candidate, err)
					}
					return
				}
			}

			parent := filepath.Dir(root)
			if parent == root {
				break
			}
			root = parent
		}
	})

	return loadErr
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

func LoadYAML[T any](path string) (*T, error) {
	data, err := ReadExpandedYAML(path)
	if err != nil {
		return nil, err
	}

	var cfg T
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ExpandedYAMLReader(path string) (*bytes.Reader, error) {
	data, err := ReadExpandedYAML(path)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}
