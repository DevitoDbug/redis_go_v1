package aof

func (a *Aof) WriteFile(value []byte) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, err := a.file.Write(value)
	if err != nil {
		return err
	}

	return nil
}
