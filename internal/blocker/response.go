package blocker

import "github.com/miekg/dns"

func BlockResponse(req *dns.Msg) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(req)

	for _, q := range req.Question {
		if q.Qtype == dns.TypeA {
			rr, _ := dns.NewRR(q.Name + " 60 IN A 0.0.0.0")
			m.Answer = append(m.Answer, rr)
		}
	}

	return m
}
