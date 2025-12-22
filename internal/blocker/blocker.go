func (b *Blocker) List() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	res := make([]string, 0, len(b.domains))
	for d := range b.domains {
		res = append(res, d)
	}
	return res
}

func (b *Blocker) Add(domain string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	domain = normalize(domain)
	b.domains[domain] = struct{}{}
	return b.flush()
}

func (b *Blocker) Delete(domain string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	domain = normalize(domain)
	delete(b.domains, domain)
	return b.flush()
}

func (b *Blocker) flush() error {
	f, err := os.Create(b.file)
	if err != nil {
		return err
	}
	defer f.Close()

	for d := range b.domains {
		_, _ = f.WriteString("0.0.0.0 " + d + "\n")
	}
	return nil
}
