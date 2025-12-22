package dnsserver

import (
	"log"
	"net/http"

	"dns-core/internal/blocker"
	"dns-core/internal/cache"
	"dns-core/internal/doh"
	"dns-core/internal/dot"
	"dns-core/internal/doq"
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

	// Web
	web.New(b).Start(":8080")

	// DoH / DoH3
	dohSrv := doh.New(r.Resolve)
	mux := http.NewServeMux()
	mux.HandleFunc("/dns-query", dohSrv.Handler)

	tlsCfg, acme := doh.NewTLSConfig(doh.TLSConfig{
		Domain: "你的域名.com", // ← 改
		Cache:  "./cert-cache",
	})

	go doh.ListenHTTPS(":443", mux, tlsCfg, acme)
	go doh.ListenHTTP3(":443", mux, tlsCfg)

	// DoT
	go func() {
		if err := dot.Listen(":853", tlsCfg, r.Resolve); err != nil {
			log.Fatal(err)
		}
	}()

	// DoQ
	go func() {
		if err := doq.Listen(":853", tlsCfg, r.Resolve); err != nil {
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
