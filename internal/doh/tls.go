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

func NewTLSConfig(cfg TLSConfig) (*tls.Config, *autocert.Manager) {
	m := &autocert.Manager{
		Cache:      autocert.DirCache(cfg.Cache),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(cfg.Domain),
	}

	tlsCfg := &tls.Config{
		GetCertificate: m.GetCertificate,
		MinVersion:     tls.VersionTLS13, // HTTP/3 必须 TLS 1.3
		NextProtos:     []string{"h3", "h3-29", "h3-30", "h3-32"},
	}

	return tlsCfg, m
}

func ListenHTTPS(addr string, handler http.Handler, tlsCfg *tls.Config, m *autocert.Manager) error {
	server := &http.Server{
		Addr:      addr,
		Handler:   handler,
		TLSConfig: tlsCfg,
	}

	go func() {
		log.Println("ACME HTTP challenge on :80")
		_ = http.ListenAndServe(":80", m.HTTPHandler(nil))
	}()

	log.Println("DoH HTTPS listening on", addr)
	return server.ListenAndServeTLS("", "")
}
