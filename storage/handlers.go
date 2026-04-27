package storage

func (s *Storage) StoreVal(key, value string) {
	s.storeLock.Lock()
	s.store[key] = value
	s.storeLock.Unlock()
}

// GetVal - returns the value from storage, if not found it will return an empty string
func (s *Storage) GetVal(key string) string {
	var val string
	s.storeLock.RLock()
	val = s.store[key]
	s.storeLock.RUnlock()

	return val
}

func (s *Storage) HStoreVal(key string, value map[string]string) {
	s.hStoreLock.Lock()
	s.hStore[key] = value
	s.hStoreLock.Unlock()
}

func (s *Storage) HGetVal(key string) map[string]string {
	var val map[string]string
	s.hStoreLock.RLock()
	val = s.hStore[key]
	s.hStoreLock.RUnlock()

	return val
}
