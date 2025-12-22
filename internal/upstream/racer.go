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

			r, _, err := up.Exchange(req)
			select {
			case ch <- Result{Resp: r, Err: err, Addr: up.Addr}:
			case <-ctx.Done():
				// 被取消：正常情况
			}
		}(u)
	}

	timeout := time.After(6 * time.Second)

	for i := 0; i < len(ups); i++ {
		select {
		case r := <-ch:
			if r.Err == nil && r.Resp != nil {
				cancel() // 主动取消其它
				return r.Resp, nil
			}
			// 记录真实错误（后面面板用）
		case <-timeout:
			return nil, errors.New("all upstreams timeout")
		}
	}

	return nil, errors.New("no valid upstream response")
}
