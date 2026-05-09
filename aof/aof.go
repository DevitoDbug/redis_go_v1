// Package aof handles data persistence
//   - stores all the executed commands into a file.
//   - stored commands are in resp format
package aof

import (
	"log"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	mu   sync.Mutex
}

func NewAof() (*Aof, error) {
	filePath := "aof/db.aof"

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: file,
		mu:   sync.Mutex{},
	}

	go func() {
		for {
			aof.mu.Lock()
			err := aof.file.Sync() // Basically make sure anything in the buffers and cache are dumped to the file
			if err != nil {
				log.Printf("failed to sync file. Error: %v", err)
			}
			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}
