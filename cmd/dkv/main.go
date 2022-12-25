package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/youngjoon-lee/dkv/service"
)

func main() {
	log.Info("starting distributed key-value store...")

	svc, err := service.New()
	if err != nil {
		log.Fatalf("failed to init service: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	sig := <-sigCh
	log.Infof("signal(%v) detected. starting graceful shutdown...", sig)

	svc.Close()
}
