package dnsserver

import (
	"log"
	"net/http"

	"dns-core/internal/blocker"
	"dns-core/internal/cache"
	"dns-core/internal/doh"
	"dns-core/internal/web"

	"github.com/miekg/dns"
)

type Server struct {
	addr     string
	resolver *Resolver
}

func New(addr string) *Server {
	b := blocker.New()
	_ = b.LoadFile("assets/block/block.txt")
	_ = b.Watch()

	c := cache.New()
	r := NewResolver(b, c)

	// Web panel
	web.New(b).Start(":8080")

	// DoH handler
	dohServer := doh.New(r.Resolve)
	mux := http.NewServeMux()
	mux.HandleFunc("/dns-query", dohServer.Handler)

	// TLS
	tlsCfg, acme := doh.NewTLSConfig(doh.TLSConfig{
		Domain: "你的域名.com", // ← 改
		Cache:  "./cert-cache",
	})

	// HTTPS (HTTP/2)
	go func() {
		if err := doh.ListenHTTPS(":443", mux, tlsCfg, acme); err != nil {
			log.Fatal(err)
		}
	}()

	// HTTP/3 (QUIC)
	go func() {
		if err := doh.ListenHTTP3(":443", mux, tlsCfg); err != nil {
			log.Fatal(err)
		}
	}()

	return &Server{
		addr:     addr,
		resolver: r,
	}
}

func (s *Server) Start() error {
	handler := dns.NewServeMux()
	handler.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		resp, err := s.resolver.Resolve(r)
		if err == nil && resp != nil {
			_ = w.WriteMsg(resp)
		}
	})

	go func() {
		log.Println("TCP DNS listening on", s.addr)
		_ = (&dns.Server{
			Addr:    s.addr,
			Net:     "tcp",
			Handler: handler,
		}).ListenAndServe()
	}()

	log.Println("UDP DNS listening on", s.addr)
	return (&dns.Server{
		Addr:    s.addr,
		Net:     "udp",
		Handler: handler,
	}).ListenAndServe()
}
