package aof

import (
	"io"
)

func (a *Aof) Read(callback func(io.Reader) error) error {
	a.mu.Lock()

	err := callback(a.file)

	a.mu.Unlock()
	return err
}
