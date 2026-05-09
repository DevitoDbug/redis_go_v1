// Package aof handles data persistence
//   - stores all the executed commands into a file.
//   - stored commands are in resp format
package aof

import (
	"bufio"
	"log"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file   *os.File
	writer *bufio.Writer
	mu     sync.Mutex
}

func NewAof() (*Aof, error) {
	filePath := "aof/db.aof"

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file:   file,
		mu:     sync.Mutex{},
		writer: bufio.NewWriter(file),
	}

	go func() {
		for {
			aof.mu.Lock()
			// Data still in the bufio buffer should be flashed into the os
			if err := aof.writer.Flush(); err != nil {
				log.Printf("failed to flush aof buffer. Error: %v", err)
			}

			// Basically make sure anything in the kernel buffers and cache are dumped to the os
			if err := aof.file.Sync(); err != nil {
				log.Printf("failed to sync file. Error: %v", err)
			}
			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}
