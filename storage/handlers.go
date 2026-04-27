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

func (s *Storage) HGetVal(hkey, hStringKey string) string {
	var valMap map[string]string
	var val string
	s.hStoreLock.RLock()
	valMap, ok := s.hStore[hkey]
	if !ok {
		return ""
	}
	val, ok = valMap[hStringKey]
	if !ok {
		return ""
	}
	s.hStoreLock.RUnlock()

	return val
}
