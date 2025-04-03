package main

import (
	"gamebooks/pkg/executor"
	"gamebooks/pkg/markdown"
	"gamebooks/pkg/repo"
	"gamebooks/pkg/storage/gorm"
	"gamebooks/pkg/web"
	"github.com/rs/zerolog"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"gorm.io/driver/sqlite"
)

func main() {
	lw := zerolog.NewConsoleWriter()
	l := zerolog.New(lw)
	l.Level(zerolog.DebugLevel)

	p := executor.New()
	r := repo.NewWithLiveReload(p)

	// database
	//s, err := bbolt.New("bbolt.db")
	//if err != nil {
	//	panic(err)
	//}
	gd := sqlite.Open("sqlite.db")
	s, err := gorm.New(gd)
	if err != nil {
		panic(err)
	}

	m := goldmark.New(
		goldmark.WithExtensions(
			meta.New(),
			extension.NewTable(),
			markdown.NewLinkTracker(r),
			extension.TaskList,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	e, err := web.New(r, s, p, l, m)
	if err != nil {
		panic(err)
	}

	if err := e.Start("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}
