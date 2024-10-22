package main

import (
	"gamebooks/pkg/executor"
	bookRepo "gamebooks/pkg/repo"
	"gamebooks/pkg/storage"
	"gamebooks/pkg/web"
)

func main() {
	p := executor.New()
	g := bookRepo.NewWithLiveReload(p)
	s := storage.NewInMemory()

	e, err := web.New(g, s, p)
	if err != nil {
		panic(err)
	}

	if err := e.Start("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}
