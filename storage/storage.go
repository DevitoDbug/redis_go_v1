// Package storage - handle the storing of the key value pairs and the
// lock management
package storage

import (
	"fmt"
	"sync"
)

type Storage struct {
	store     map[string]string
	storeLock sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		store:     make(map[string]string),
		storeLock: sync.RWMutex{},
	}
}

func (s *Storage) PrintStore() {
	fmt.Printf("\n*******************\nstore:\n%v\n*******************\n", s.store)
}
