package blocker

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

type Blocker struct {
	mu     sync.RWMutex
	domains map[string]struct{}
}

func New() *Blocker {
	return &Blocker{
		domains: make(map[string]struct{}),
	}
}

func (b *Blocker) LoadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	tmp := make(map[string]struct{})
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			domain := strings.TrimSuffix(fields[1], ".")
			tmp[domain] = struct{}{}
		}
	}

	b.mu.Lock()
	b.domains = tmp
	b.mu.Unlock()

	return nil
}

func (b *Blocker) IsBlocked(domain string) bool {
	domain = strings.TrimSuffix(domain, ".")
	b.mu.RLock()
	_, ok := b.domains[domain]
	b.mu.RUnlock()
	return ok
}
