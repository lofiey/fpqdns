package ecs

import "github.com/miekg/dns"

func Attach(msg *dns.Msg, ip string, mask uint8) {
	opt := msg.IsEdns0()
	if opt == nil {
		opt = &dns.OPT{}
		opt.Hdr.Name = "."
		opt.Hdr.Rrtype = dns.TypeOPT
		msg.Extra = append(msg.Extra, opt)
	}

	subnet := &dns.EDNS0_SUBNET{
		Code:          dns.EDNS0SUBNET,
		Address:       ip,
		SourceNetmask: mask,
	}
	opt.Option = append(opt.Option, subnet)
}
