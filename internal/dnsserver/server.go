package dnsserver

import (
	"log"

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

	web.New(b).Start(":8080")

	// DoH
	dohServer := doh.New(r.Resolve)
	go func() {
		log.Println("DoH listening on :443 (/dns-query)")
		http.HandleFunc("/dns-query", dohServer.Handler)
		_ = http.ListenAndServe(":443", nil)
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
