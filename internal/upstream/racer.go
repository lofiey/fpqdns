package upstream

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/miekg/dns"
)

type Result struct {
	Resp *dns.Msg
	Err  error
	Addr string
}

func Race(req *dns.Msg, ups []Upstream) (*dns.Msg, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan Result, len(ups))
	var wg sync.WaitGroup

	for _, u := range ups {
		wg.Add(1)
		go func(up Upstream) {
			defer wg.Done()

			resp, _, err := up.Exchange(req)
			select {
			case ch <- Result{Resp: resp, Err: err, Addr: up.Addr}:
			case <-ctx.Done():
				// 主动取消，不算错误
			}
		}(u)
	}

	timeout := time.After(6 * time.Second)

	for i := 0; i < len(ups); i++ {
		select {
		case r := <-ch:
			if r.Err == nil && r.Resp != nil {
				cancel()
				return r.Resp, nil
			}
		case <-timeout:
			return nil, errors.New("all upstreams timeout")
		}
	}

	return nil, errors.New("no valid upstream response")
}
