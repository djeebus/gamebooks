package storage

import (
	"github.com/rs/zerolog/log"
)

func GetString(s Storage, key string) string {
	result, err := s.Get(key)
	if err != nil {
		log.Error().Err(err).Msg("failed to get string")
	}
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
	result, err := s.Get(key)
	if err != nil {
		log.Error().Err(err).Msg("failed to get bool")
	}

	if result == nil {
		return false
	}

	val, ok := result.(bool)
	if !ok {
		log.Error().Err(err).Msg("value was not a bool")
		return false
	}

	return val
}

func GetSlice[T any](s Storage, key string) []T {
	val, err := s.Get(key)
	if err != nil {
		log.Error().Err(err).Msg("failed to get value")
	}
	items, ok := val.([]T)
	if !ok || items == nil {
		items = make([]T, 0)
	}
	return items
}

func Push[T any](s Storage, key string, value T) error {
	items := GetSlice[T](s, key)
	items = append(items, value)
	return s.Set(key, items)
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

func Pop[T any](s Storage, key string) (T, error) {
	items := GetSlice[T](s, key)
	if len(items) == 0 {
		var t T
		return t, nil
	}

	popped := items[len(items)-1]
	items = items[:len(items)-1]
	if err := s.Set(key, items); err != nil {
		return popped, err
	}

	return popped, nil
}
