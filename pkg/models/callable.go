package models

type Callable func(args []any, kwargs map[string]any) (any, error)
