package cache

import "github.com/miekg/dns"

// 尝试从已有缓存中拼出 CNAME 链
func ResolveCNAME(q dns.Question, c *Cache) *dns.Msg {
	msg := new(dns.Msg)
	msg.SetQuestion(q.Name, q.Qtype)

	current := q.Name

	for i := 0; i < 8; i++ { // 防止死循环
		cached := c.Get(dns.Question{
			Name:   current,
			Qtype: dns.TypeCNAME,
			Qclass: dns.ClassINET,
		})
		if cached == nil || len(cached.Answer) == 0 {
			break
		}

		msg.Answer = append(msg.Answer, cached.Answer...)

		cname, ok := cached.Answer[0].(*dns.CNAME)
		if !ok {
			break
		}
		current = cname.Target
	}

	final := c.Get(dns.Question{
		Name:   current,
		Qtype: q.Qtype,
		Qclass: dns.ClassINET,
	})
	if final != nil {
		msg.Answer = append(msg.Answer, final.Answer...)
		return msg
	}

	return nil
}
