package executor

import (
	"fmt"
	"gamebooks/pkg/storage"
	"github.com/pkg/errors"
	"go.starlark.net/starlark"
)

func starlarkPredeclared(s storage.Storage) starlark.StringDict {
	return starlark.StringDict{
		"dice_roll": starlark.NewBuiltin("dice_roll", func(t *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var count, size int

			if err := starlark.UnpackArgs(b.Name(),
				args, kwargs,
				"count", &count,
				"size", &size,
			); err != nil {
				return nil, errors.Wrap(err, "failed to parse diceRoll args")
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
				return nil, errors.Wrap(err, "failed to parse storageSet args")
			}

			value := s.Get(key)
			return makeStarlarkValue(value)
		}),
		"storage_set": starlark.NewBuiltin("storage_set", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var key string
			var value starlark.Value

			if err := starlark.UnpackArgs(fn.Name(),
				args, kwargs,
				"key", &key,
				"value", &value,
			); err != nil {
				return nil, errors.Wrap(err, "failed to parse storageSet args")
			}

			val, err := unwrapStarlarkValue(value)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unwrap starlark value")
			}

			s.Set(key, val)
			return starlark.None, nil
		}),
	}
}

func unwrapStarlarkValue(value starlark.Value) (interface{}, error) {
	switch t := value.(type) {
	case starlark.Bool:
		return bool(t), nil
	case starlark.Int:
		val, ok := t.Int64()
		if !ok {
			return nil, errors.New("int64 overflow")
		}
		return val, nil
	case starlark.String:
		return string(t), nil
	default:
		return nil, fmt.Errorf("unknown type: %v", value)
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
	}

	return nil, fmt.Errorf("unknown type: %v", value)
}
