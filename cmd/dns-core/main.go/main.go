package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"dns-core/internal/blocker"
	"dns-core/internal/cache"
	"dns-core/internal/dnsserver"
)

func main() {
	log.Println("DNS Core starting...")

	// =========================
	// Blocker
	// =========================
	block := blocker.New()

	// block 规则文件（你已有）
	if err := block.LoadFile("assets/block/block.txt"); err != nil {
		log.Println("block file load error:", err)
	}

	// 热更新监听
	if err := block.Watch(); err != nil {
		log.Println("block watch error:", err)
	}

	// =========================
	// Cache
	// =========================
	c := cache.New()

	// =========================
	// DNS Server
	// =========================
	server := dnsserver.New(":53")

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal("dns server error:", err)
		}
	}()

	log.Println("DNS Core started successfully")

	// =========================
	// Graceful Shutdown
	// =========================
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	log.Println("DNS Core shutting down...")
}
