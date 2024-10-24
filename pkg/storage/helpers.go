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

func GetSlice[T any](s Storage, key string) []T {
	val := s.Get(key)
	items, ok := val.([]T)
	if !ok || items == nil {
		items = make([]T, 0)
	}
	return items
}

func Push[T any](s Storage, key string, value T) {
	items := GetSlice[T](s, key)
	items = append(items, value)
	s.Set(key, items)
}

func Peek[T any](s Storage, key string) T {
	items := GetSlice[T](s, key)
	if len(items) == 0 {
		var t T
		return t
	}

	popped := items[len(items)-1]
	return popped
}

func Pop[T any](s Storage, key string) T {
	items := GetSlice[T](s, key)
	if len(items) == 0 {
		var t T
		return t
	}

	popped := items[len(items)-1]
	items = items[:len(items)-1]
	s.Set(key, items)
	return popped
}
