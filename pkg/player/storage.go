package player

import (
	"gamebooks/pkg/storage"
	"github.com/Shopify/go-lua"
)

func getStorage(s storage.Storage) lua.Function {
	return func(l *lua.State) int {
		key := lua.CheckString(l, 1)
		value := s.Get(key)
		return deepPush(l, value)
	}
}

func setStorage(s storage.Storage) lua.Function {
	return func(l *lua.State) int {
		key := lua.CheckString(l, 1)
		val := l.ToValue(2)
		s.Set(key, val)
		return 1
	}
}
