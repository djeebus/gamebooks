package executor

import (
	"fmt"
	"gamebooks/pkg/models"
	"gamebooks/pkg/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.starlark.net/starlark"
)

var ErrCannotCallPageFunction = errors.New("cannot call page function in this context")

func starlarkPredeclared(s storage.Storage, page *models.Page) starlark.StringDict {
	var pageStorage storage.Storage
	if page != nil {
		pageStorage = storage.NamespacedStorage(s, page.PageID)
	}

	return starlark.StringDict{
		"dice_roll": starlark.NewBuiltin("dice_roll", func(t *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var count, size int

			if err := starlark.UnpackArgs(b.Name(),
				args, kwargs,
				"count", &count,
				"size", &size,
			); err != nil {
				return nil, errors.Wrap(err, "failed to parse dice_roll args")
			}

			total := rollDie(count, size)

			return starlark.MakeInt(total), nil
		}),
		"storage_get": starlark.NewBuiltin("storage_get", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var key string

			if err := starlark.UnpackArgs(fn.Name(),
				args, kwargs,
				"key", &key,
			); err != nil {
				return nil, errors.Wrap(err, "failed to parse storage_get args")
			}

			value, err := s.Get(key)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get %q from storage", key)
			}
			return makeStarlarkValue(value)
		}),
		"storage_page_get": starlark.NewBuiltin("storage_page_get", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var key string

			if err := starlark.UnpackArgs(fn.Name(),
				args, kwargs,
				"key", &key,
			); err != nil {
				return nil, errors.Wrap(err, "failed to parse storage_get args")
			}

			if pageStorage == nil {
				return nil, ErrCannotCallPageFunction
			}

			value, err := pageStorage.Get(key)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get value")
			}
			return makeStarlarkValue(value)
		}),
		"storage_page_remove": starlark.NewBuiltin("storage_page_remove", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var key string

			if err := starlark.UnpackArgs(fn.Name(),
				args, kwargs,
				"key", &key,
			); err != nil {
				return nil, errors.Wrap(err, "failed to parse storage_page_remove args")
			}

			if pageStorage == nil {
				return nil, ErrCannotCallPageFunction
			}

			if err := pageStorage.Remove(key); err != nil {
				return nil, errors.Wrap(err, "failed to remove key")
			}

			return starlark.None, nil
		}),
		"storage_page_set": starlark.NewBuiltin("storage_page_set", func(t *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var key string
			var value starlark.Value

			if err := starlark.UnpackArgs(fn.Name(),
				args, kwargs,
				"key", &key,
				"value", &value,
			); err != nil {
				return nil, errors.Wrap(err, "failed to parse storage_page_set args")
			}

			val, err := unwrapStarlarkValue(t, value)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unwrap starlark value")
			}

			if pageStorage == nil {
				return nil, ErrCannotCallPageFunction
			}

			if err = pageStorage.Set(key, val); err != nil {
				return nil, errors.Wrap(err, "failed to set key")
			}
			return starlark.None, nil
		}),
		"storage_remove": starlark.NewBuiltin("storage_remove", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var key string

			if err := starlark.UnpackArgs(fn.Name(),
				args, kwargs,
				"key", &key,
			); err != nil {
				return nil, errors.Wrap(err, "failed to parse storage_remove args")
			}

			if err := s.Remove(key); err != nil {
				return nil, errors.Wrap(err, "failed to remove key")
			}

			return starlark.None, nil
		}),
		"storage_pop": starlark.NewBuiltin("storage_pop", func(t *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var key string

			if err := starlark.UnpackArgs(fn.Name(),
				args, kwargs,
				"key", &key,
			); err != nil {
				return nil, errors.Wrap(err, "failed to parse storage_set args")
			}

			if _, err := storage.Pop[any](s, key); err != nil {
				return nil, errors.Wrap(err, "failed to set key")
			}
			return starlark.None, nil
		}),
		"storage_push": starlark.NewBuiltin("storage_push", func(t *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var key string
			var value starlark.Value

			if err := starlark.UnpackArgs(fn.Name(),
				args, kwargs,
				"key", &key,
				"value", &value,
			); err != nil {
				return nil, errors.Wrap(err, "failed to parse storage_set args")
			}

			val, err := unwrapStarlarkValue(t, value)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unwrap starlark value")
			}

			if err = storage.Push(s, key, val); err != nil {
				return nil, errors.Wrap(err, "failed to set key")
			}
			return starlark.None, nil
		}),
		"storage_set": starlark.NewBuiltin("storage_set", func(t *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var key string
			var value starlark.Value

			if err := starlark.UnpackArgs(fn.Name(),
				args, kwargs,
				"key", &key,
				"value", &value,
			); err != nil {
				return nil, errors.Wrap(err, "failed to parse storage_set args")
			}

			val, err := unwrapStarlarkValue(t, value)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unwrap starlark value")
			}

			if err = s.Set(key, val); err != nil {
				return nil, errors.Wrap(err, "failed to set key")
			}
			return starlark.None, nil
		}),
		"debug": starlark.NewBuiltin("debug", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var message string

			if err := starlark.UnpackArgs(fn.Name(),
				args, kwargs,
				"message", &message,
			); err != nil {
				return nil, errors.Wrap(err, "failed to parse debug args")
			}

			log.Debug().Msg(message)
			return starlark.None, nil
		}),
	}
}

func unwrapStarlarkDict(t *starlark.Thread, d starlark.StringDict) (map[string]any, error) {
	var err error
	result := map[string]any{}
	for key, value := range d {
		if result[key], err = unwrapStarlarkValue(t, value); err != nil {
			return nil, errors.Wrapf(err, "failed to convert %s", key)
		}
	}
	return result, nil
}

func unwrapStarlarkValue(t *starlark.Thread, value starlark.Value) (interface{}, error) {
	switch v := value.(type) {
	case starlark.Bool:
		return bool(v), nil
	case *starlark.Dict:
		d := make(map[string]interface{})
		for _, key := range v.Keys() {
			val, ok, err := v.Get(key)
			if !ok {
				continue
			}
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get %v from dict", key)
			}

			keystr, ok := key.(starlark.String)
			if !ok {
				log.Warn().Any("key", val).Msg("non-string key")
				continue
			}

			unwrappedValue, err := unwrapStarlarkValue(t, val)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to unwrap %v from dict", key)
			}

			d[string(keystr)] = unwrappedValue
		}
		return d, nil
	case *starlark.Function:
		return newCallable(t, v), nil
	case starlark.Int:
		val, ok := v.Int64()
		if !ok {
			return nil, errors.New("int64 overflow")
		}
		return val, nil
	case *starlark.List:
		var items []any
		for item := range v.Elements() {
			unwrapped, err := unwrapStarlarkValue(t, item)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unwrap item")
			}
			items = append(items, unwrapped)
		}
		return items, nil
	case starlark.NoneType:
		return nil, nil
	case starlark.String:
		return string(v), nil
	default:
		return nil, fmt.Errorf("cannot convert %T to golang value", value)
	}
}

func newCallable(t *starlark.Thread, fn starlark.Callable) models.Callable {
	return func(args []any, kwargs map[string]any) (any, error) {
		sargs, err := makeStarlarkValueList(args)
		if err != nil {
			return nil, errors.Wrap(err, "failed to wrap args")
		}

		skwargs, err := makeStarlarkKwargsTupleSlice(kwargs)
		if err != nil {
			return nil, errors.Wrap(err, "failed to wrap kwargs")
		}

		result, err := starlark.Call(t, fn, sargs, skwargs)
		if err != nil {
			return nil, err
		}

		return unwrapStarlarkValue(t, result)
	}
}

func makeStarlarkValue(value interface{}) (starlark.Value, error) {
	if value == nil {
		return starlark.None, nil
	}

	switch v := value.(type) {
	case bool:
		return starlark.Bool(v), nil
	case string:
		return starlark.String(v), nil
	case int:
		return starlark.MakeInt(v), nil
	case int64:
		return starlark.MakeInt64(v), nil
	case float64:
		return starlark.Float(v), nil
	case map[string]interface{}:
		var result starlark.Dict
		for key, value := range v {
			skey, err := makeStarlarkValue(key)
			if err != nil {
				return nil, errors.Wrap(err, "failed to wrap key")
			}

			svalue, err := makeStarlarkValue(value)
			if err != nil {
				return nil, errors.Wrap(err, "failed to wrap value")
			}

			if err = result.SetKey(skey, svalue); err != nil {
				return nil, errors.Wrapf(err, "failed to set key %q", key)
			}
		}
		return &result, nil
	case []any:
		return makeStarlarkValueList(v)
	}

	return nil, fmt.Errorf("cannot convert %T to starlark value", value)
}

func makeStarlarkKwargsTupleSlice(value map[string]any) ([]starlark.Tuple, error) {
	var slice []starlark.Tuple
	for key, value := range value {
		keywrap, err := makeStarlarkValue(key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to wrap key")
		}

		valuewrap, err := makeStarlarkValue(value)
		if err != nil {
			return nil, errors.Wrap(err, "failed to wrap value")
		}

		slice = append(slice, []starlark.Value{keywrap, valuewrap})
	}
	return slice, nil
}

func makeStarlarkValueList(value []any) (starlark.Tuple, error) {
	var results []starlark.Value
	for _, item := range value {
		wrapped, err := makeStarlarkValue(item)
		if err != nil {
			return nil, errors.Wrap(err, "failed to make value")
		}
		results = append(results, wrapped)
	}
	return results, nil
}
