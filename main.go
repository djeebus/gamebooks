package main

import (
	bookRepo "gamebooks/pkg/repo"
	"gamebooks/pkg/storage"
	"gamebooks/pkg/web"
)

func main() {
	g := bookRepo.NewWithLiveReload()
	s := storage.NewInMemory()

	e, err := web.New(g, s)
	if err != nil {
		panic(err)
	}

	if err := e.Start("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}
