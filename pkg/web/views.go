package web

import (
	"embed"
	"gamebooks/pkg/executor"
	bookRepo "gamebooks/pkg/repo"
	"gamebooks/pkg/storage"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
)

//go:embed templates/*
var fs embed.FS

type views struct {
	game     bookRepo.Game
	storage  storage.Storage
	player   *executor.Player
	markdown goldmark.Markdown
}

func newViews(game bookRepo.Game, storage storage.Storage, player *executor.Player) (*views, error) {
	return &views{
		game:    game,
		storage: storage,
		player:  player,
		markdown: goldmark.New(
			goldmark.WithExtensions(
				extension.NewTable(),
				meta.Meta,
				NewLinkTracker(),
			),
		),
	}, nil
}
