package aof

import "os"

func (a *Aof) Read(callback func(*os.File) error) error {
	a.mu.Lock()

	err := callback(a.file)

	a.mu.Unlock()
	return err
}
