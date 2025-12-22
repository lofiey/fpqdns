package dnsserver

import (
	"log"

	"github.com/miekg/dns"
)

func (s *Server) listenDoT(certFile, keyFile string) {
	handler := dns.NewServeMux()
	handler.HandleFunc(".", s.handleDNS)

	server := &dns.Server{
		Addr:      ":853",
		Net:       "tcp-tls",
		Handler:   handler,
		TLSConfig: s.tlsConfig(certFile, keyFile),
	}

	log.Println("DoT listening on :853")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("DoT failed:", err)
	}
}
