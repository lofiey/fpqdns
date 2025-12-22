package app

import (
	"log"

	"fpqdns/internal/config"
	"fpqdns/internal/dnsserver"
)

func Run() error {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		return err
	}

	srv := dnsserver.New(cfg.Listen.DNS)
	if err := srv.Start(); err != nil {
		return err
	}

	log.Println("[app] fpqdns started")
	select {}
}
