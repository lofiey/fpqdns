package dnsserver

import (
	"log"
)

type Server struct {
	addr string
}

func New(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) Start() error {
	go s.listenUDP()
	go s.listenTCP()
	go s.listenDoH()

	log.Println("[dns] server started on", s.addr)
	return nil
}
