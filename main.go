package main

import (
	"gamebooks/pkg/executor"
	bookRepo "gamebooks/pkg/repo"
	"gamebooks/pkg/storage"
	"gamebooks/pkg/web"
	zerolog "github.com/rs/zerolog"
)

func main() {
	lw := zerolog.NewConsoleWriter()
	l := zerolog.New(lw)
	l.Level(zerolog.DebugLevel)

	p := executor.New()
	g := bookRepo.NewWithLiveReload(p)
	s := storage.NewInMemory()

	e, err := web.New(g, s, p, l)
	if err != nil {
		panic(err)
	}

	if err := e.Start("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}
