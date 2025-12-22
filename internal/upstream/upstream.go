package upstream

import (
	"time"

	"github.com/miekg/dns"
)

type Upstream struct {
	Addr string // 1.1.1.1:53 / tls:// / https://
	Net  string // udp / tcp
}

func (u *Upstream) Exchange(req *dns.Msg) (*dns.Msg, error) {
	c := &dns.Client{
		Net:     u.Net,
		Timeout: 5 * time.Second,
	}
	return c.Exchange(req, u.Addr)
}
