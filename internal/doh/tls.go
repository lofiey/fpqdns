package doh

import (
	"crypto/tls"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

type TLSConfig struct {
	Domain string
	Cache  string
}

func ListenHTTPS(addr string, handler http.Handler, cfg TLSConfig) error {
	m := &autocert.Manager{
		Cache:      autocert.DirCache(cfg.Cache),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(cfg.Domain),
	}

	server := &http.Server{
		Addr: addr,
		TLSConfig: &tls.Config{
			GetCertificate: m.GetCertificate,
			MinVersion:     tls.VersionTLS12,
		},
		Handler: handler,
	}

	// HTTP-01 challenge
	go func() {
		log.Println("ACME HTTP challenge on :80")
		_ = http.ListenAndServe(":80", m.HTTPHandler(nil))
	}()

	log.Println("DoH HTTPS listening on", addr)
	return server.ListenAndServeTLS("", "")
}
