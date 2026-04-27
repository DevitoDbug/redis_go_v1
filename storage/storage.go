// Package storage - handle the storing of the key value pairs and the
// lock management
package storage

import "sync"

type Storage struct {
	Store     map[string]string
	StoreLock sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		Store:     make(map[string]string),
		StoreLock: sync.RWMutex{},
	}
}
