package doq

import (
	"context"
	"log"

	"github.com/miekg/dns"
	"github.com/quic-go/quic-go"
)

type Resolver func(*dns.Msg) (*dns.Msg, error)

func Listen(addr string, tlsCfg *tls.Config, r Resolver) error {
	listener, err := quic.ListenAddr(addr, tlsCfg, &quic.Config{
		EnableDatagrams: true,
	})
	if err != nil {
		return err
	}

	log.Println("DoQ listening on", addr)

	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			continue
		}

		go func() {
			for {
				stream, err := sess.AcceptStream(context.Background())
				if err != nil {
					return
				}

				go handleStream(stream, r)
			}
		}()
	}
}

func handleStream(s quic.Stream, r Resolver) {
	defer s.Close()

	buf := make([]byte, 4096)
	n, err := s.Read(buf)
	if err != nil {
		return
	}

	req := new(dns.Msg)
	if err := req.Unpack(buf[:n]); err != nil {
		return
	}

	resp, err := r(req)
	if err != nil {
		return
	}

	out, err := resp.Pack()
	if err != nil {
		return
	}

	_, _ = s.Write(out)
}
