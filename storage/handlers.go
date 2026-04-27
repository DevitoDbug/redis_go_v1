package storage

func (s *Storage) StoreVal(key, value string) {
	s.storeLock.Lock()
	defer s.storeLock.Unlock()

	s.store[key] = value
}

// GetVal - returns the value from storage, if not found it will return an empty string
func (s *Storage) GetVal(key string) string {
	var val string
	s.storeLock.RLock()
	defer s.storeLock.RUnlock()

	val = s.store[key]
	return val
}

func (s *Storage) HStoreVal(hKey, hStringKey, value string) {
	s.hStoreLock.Lock()
	defer s.hStoreLock.Unlock()

	innerMap, ok := s.hStore[hKey]
	if !ok || innerMap == nil {
		innerMap = make(map[string]string)
		s.hStore[hKey] = innerMap
	}

	innerMap[hStringKey] = value
}

func (s *Storage) HGetVal(hkey, hStringKey string) string {
	var valMap map[string]string
	var val string
	s.hStoreLock.RLock()
	defer s.hStoreLock.RUnlock()

	valMap, ok := s.hStore[hkey]
	if !ok {
		return ""
	}
	val, ok = valMap[hStringKey]
	if !ok {
		return ""
	}

	return val
}

func (s *Storage) HGetAllVal(hkey string) map[string]string {
	var val map[string]string
	s.hStoreLock.RLock()
	defer s.hStoreLock.RUnlock()

	val, ok := s.hStore[hkey]
	if !ok {
		return nil
	}

	return val
}
