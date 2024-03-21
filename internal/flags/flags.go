package flags

import (
	"flag"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
	"time"
)

// ConfigVariables содержит конфигурационные переменные приложения.
type ConfigVariables struct {
	AddressEnv        string `env:"ADDRESS"`
	Address           string
	ReportIntervalEnv int `env:"REPORT_INTERVAL"`
	ReportInterval    int
	PollIntervalEnv   int `env:"POLL_INTERVAL"`
	PollInterval      int
}

// ParseFlags парсит флаги и переменные окружения для настройки конфигурации.
func ParseFlags(typeSvc string) (string, time.Duration, time.Duration) {
	var cfg ConfigVariables
	urlSchema := "http://"
	defaultAddress := "localhost:8080"

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	if typeSvc == "server" {
		defaultAddress = ":8080"
	}
	flag.StringVar(&cfg.Address, "a", defaultAddress, "HTTP server address")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "Report interval")
	flag.IntVar(&cfg.PollInterval, "p", 2, "Poll interval")
	flag.Parse()
	if cfg.AddressEnv != "" {
		cfg.Address = cfg.AddressEnv
	}
	if typeSvc == "agent" {
		cfg.Address = urlSchema + cfg.Address
	}

	if cfg.ReportIntervalEnv != 0 {
		cfg.ReportInterval = cfg.ReportIntervalEnv
	}

	if cfg.PollIntervalEnv != 0 {
		cfg.PollInterval = cfg.PollIntervalEnv
	}
	return cfg.Address, time.Duration(cfg.ReportInterval) * time.Second, time.Duration(cfg.PollInterval) * time.Second
}
