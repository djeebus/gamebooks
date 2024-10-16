package main

import (
	"gamebooks/pkg/game"
	"gamebooks/pkg/web"
	"github.com/pkg/errors"
)

func main() {
	g, err := game.New(false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create game")
	}

	e, err := web.New(g)
	if err != nil {
		panic(err)
	}

	if err := e.Start("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}
