package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type LogLevel log.Level

type Config struct {
	LogLevel   LogLevel `envconfig:"LOG_LEVEL" default:"info"`
	RPCPort    int      `envconfig:"RPC_PORT" required:"true"`
	RESTPort   int      `envconfig:"REST_PORT" required:"true"`
	DBPath     string   `envconfig:"DB_PATH" required:"true"`
	LeaderAddr string   `envconfig:"LEADER_ADDR"`
}

func (l *LogLevel) Decode(value string) error {
	level, err := log.ParseLevel(value)
	if err != nil {
		return fmt.Errorf("failed to parse LogLevel: %w", err)
	}

	*l = LogLevel(level)
	return nil
}
