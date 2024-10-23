package executor

import (
	"gamebooks/pkg/storage"
	"github.com/Shopify/go-lua"
)

func storageFunctions(s storage.Storage) []lua.RegistryFunction {
	return []lua.RegistryFunction{
		{Name: "get", Function: getStorage(s)},
		{Name: "set", Function: setStorage(s)},
	}
}

func pageStorageFunctions(s storage.Storage, pageID string) []lua.RegistryFunction {
	return append([]lua.RegistryFunction{
		{Name: "get_page", Function: getStorage(storage.NamespacedStorage(s, pageID))},
		{Name: "set_page", Function: setStorage(storage.NamespacedStorage(s, pageID))},
	}, storageFunctions(s)...)
}

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
