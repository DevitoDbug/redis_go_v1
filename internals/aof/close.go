package aof

func (a *Aof) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.file.Close()
}
