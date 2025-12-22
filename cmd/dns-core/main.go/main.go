package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"dns-core/internal/dnsserver"
)

func main() {
	srv := dnsserver.New(":53")

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("dns start failed: %v", err)
		}
	}()

	log.Println("DNS server started on :53")

	// 等待退出信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("shutting down DNS server")
	srv.Stop()
}
