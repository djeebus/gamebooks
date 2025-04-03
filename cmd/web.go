package cmd

import (
	"gamebooks/pkg/container"
	"gamebooks/pkg/web"
	"github.com/rs/zerolog"
)

func http(ctr container.Container, l zerolog.Logger) {
	e, err := web.New(ctr, l)
	if err != nil {
		panic(err)
	}

	if err := e.Start("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}
