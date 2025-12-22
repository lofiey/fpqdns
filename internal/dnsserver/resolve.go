import "dns-core/internal/upstream"

// ...

if isCN {
	up := upstream.Upstream{
		Addr: "223.5.5.5:53",
		Net:  "udp",
	}
	resp, _, err := up.Exchange(req)
	if err != nil {
		return nil, err
	}
	r.cache.Set(q, resp)
	return resp, nil
}

// 海外：并发竞速
ups := []upstream.Upstream{
	{Addr: "8.8.8.8:53", Net: "udp"},
	{Addr: "1.1.1.1:53", Net: "udp"},
	{Addr: "9.9.9.9:53", Net: "udp"},
}

resp, err := upstream.Race(req, ups)
if err != nil {
	return nil, err
}

r.cache.Set(q, resp)
return resp, nil
