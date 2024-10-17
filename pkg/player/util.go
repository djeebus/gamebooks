package player

import (
	"fmt"
	"github.com/Shopify/go-lua"
)

func open(l *lua.State, libraryName string, fns []lua.RegistryFunction) {
	libraryOpen := func(l *lua.State) int {
		lua.NewLibrary(l, fns)
		return 1
	}
	lua.Require(l, libraryName, libraryOpen, false)
	l.Pop(1)
}

func deepPush(l *lua.State, v interface{}) int {
	switch val := v.(type) {
	case nil:
		l.PushNil()
	case string:
		l.PushString(val)
	case int:
		l.PushInteger(val)
	case int32:
		l.PushInteger(int(val))
	case int64:
		l.PushInteger(int(val))
	case float32:
		l.PushNumber(float64(val))
	case float64:
		l.PushNumber(val)
	default:
		panic(fmt.Sprintf("%T: %v", val, val))
	}

	return 1
}
