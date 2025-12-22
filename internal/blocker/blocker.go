package blocker

import (
	"bufio"
	"os"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

type Blocker struct {
	mu    sync.RWMutex
	rules map[string]struct{}
	file  string
}

func New() *Blocker {
	return &Blocker{
		rules: make(map[string]struct{}),
	}
}

func (b *Blocker) LoadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	rules := make(map[string]struct{})

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			domain := dns.Fqdn(fields[1])
			rules[domain] = struct{}{}
		}
	}

	b.mu.Lock()
	b.rules = rules
	b.file = path
	b.mu.Unlock()

	return nil
}

func (b *Blocker) IsBlocked(name string) bool {
	name = dns.Fqdn(name)

	b.mu.RLock()
	_, ok := b.rules[name]
	b.mu.RUnlock()

	return ok
}

func (b *Blocker) List() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	out := make([]string, 0, len(b.rules))
	for k := range b.rules {
		out = append(out, k)
	}
	return out
}

func (b *Blocker) Add(domain string) error {
	domain = dns.Fqdn(domain)

	f, err := os.OpenFile(b.file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("0.0.0.0 " + domain + "\n")
	return err
}

func (b *Blocker) Remove(domain string) error {
	domain = dns.Fqdn(domain)

	file, err := os.ReadFile(b.file)
	if err != nil {
		return err
	}

	lines := strings.Split(string(file), "\n")
	var out []string
	for _, line := range lines {
		if !strings.Contains(line, domain) {
			out = append(out, line)
		}
	}

	return os.WriteFile(b.file, []byte(strings.Join(out, "\n")), 0644)
}

func BlockResponse(req *dns.Msg) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(req)
	m.Authoritative = true
	return m
}
