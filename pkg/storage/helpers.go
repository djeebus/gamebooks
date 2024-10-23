package storage

func GetString(s Storage, key string) string {
	result := s.Get(key)
	if result == nil {
		return ""
	}

	val, ok := result.(string)
	if !ok {
		return ""
	}

	return val
}
