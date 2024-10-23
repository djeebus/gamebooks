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

func GetBool(s Storage, key string) bool {
	result := s.Get(key)
	if result == nil {
		return false
	}

	val, ok := result.(bool)
	if !ok {
		return false
	}

	return val

}
