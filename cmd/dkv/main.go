package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/app"
	"github.com/youngjoon-lee/dkv/config"
)

const (
	envPrefix = "DKV"
)

func main() {
	log.Info("starting distributed key-value store...")

	var conf config.Config
	if err := envconfig.Process(envPrefix, &conf); err != nil {
		log.Fatalf("failed to parse env vars: %v", err)
	}

	initLogger(log.Level(conf.LogLevel))
	log.Debugf("config: %v", conf)

	app, err := app.New(conf)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}
	defer app.Close()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	sig := <-sigCh
	log.Infof("signal(%v) detected. starting graceful shutdown...", sig)
}

func initLogger(level log.Level) {
	log.SetLevel(level)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
}
