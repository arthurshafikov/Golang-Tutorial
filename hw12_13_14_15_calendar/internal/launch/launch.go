package launch

import (
	"flag"

	"github.com/arthurshafikov/hw12_13_14_15_calendar/internal/config"
	"github.com/arthurshafikov/hw12_13_14_15_calendar/internal/logger"
)

var (
	configFilePath string
	logFilePath    string
)

func init() {
	flag.StringVar(&configFilePath, "config", "./configs/config.toml", "Path to configuration file")
	flag.StringVar(&logFilePath, "log", "./logs/log.txt", "Path to log file")
}

func InitializateConfigAndLoggerFromFlags() (config.Config, *logger.Logger) {
	flag.Parse()

	config := config.NewConfig(configFilePath)
	logg := logger.New(config.Logger.Level, logFilePath)

	return config, logg
}
