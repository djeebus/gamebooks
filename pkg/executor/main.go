package executor

import (
	"github.com/pkg/errors"
)

func New() *Executor {
	return &Executor{}
}

type Executor struct {
}

var ErrNoField = errors.New("no field")
