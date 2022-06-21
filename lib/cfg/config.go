package cfg

import (
	"fmt"
	"sync"

	"cp/lib/zlog"
	"github.com/rs/zerolog"
)

type Config struct {
	Service string `env:"-"`

	// Orchestrator connection settings
	OrcHost           string `env:"ORC_HOST" envDefault:"192.168.0.50"`
	OrcPort           string `env:"ORC_PORT" envDefault:"9010"`
	OrcTlsSrvOverride string `env:"ORC_TLS_OVERRIDE" envDefault:""`

	// Listener settings
	LisHost string `env:"LISTEN_HOST"`
	LisPort string `env:"LISTEN_PORT"`

	// Db settings
	SqlType string `env:"SQL_TYPE" envDefault:"postgres"`
	SqlHost string `env:"SQL_HOST" envDefault:"192.168.0.50"`
	SqlPort string `env:"SQL_PORT" envDefault:"5432"`
	SqlUser string `env:"SQL_USER" envDefault:"postgres"`
	SqlPass string `env:"SQL_PASS" envDefault:"postgres"`
	SqlDb   string `env:"SQL_DB"`

	// Develop settings
	IsDev    bool `env:"IS_DEV" envDefault:"false"`
	IsDocker bool `env:"IS_DOCKER" envDefault:"false"`

	// Other settings
	Internal map[string]string `env:"-"`
	mu       sync.Mutex        `env:"-"`
}

func (cfg *Config) GetSrvLogger() zerolog.Logger {
	return zlog.InitLogger(cfg.IsDev).With().Str("service", cfg.Service).Logger()
}

func (cfg *Config) GetListenerAddress() string {

	if len(cfg.LisPort) == 0 {
		return ""
	}

	if cfg.IsDocker {
		return ":" + cfg.LisPort
	} else {
		return cfg.LisHost + ":" + cfg.LisPort
	}
}

func (cfg *Config) GetOrcTarget() string {
	return fmt.Sprintf("%s:%s", cfg.OrcHost, cfg.OrcPort)
}

func (cfg *Config) GetInternalSettings() map[string]string {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	settings := make(map[string]string)
	for key, value := range cfg.Internal {
		settings[key] = value
	}

	return settings
}

func (cfg *Config) SetInternalSettings(key, val string) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if cfg.Internal == nil {
		cfg.Internal = make(map[string]string)
	}

	cfg.Internal[key] = val
}
