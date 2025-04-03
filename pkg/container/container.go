package container

import (
	"gamebooks/pkg/executor"
	"gamebooks/pkg/markdown"
	"gamebooks/pkg/repo"
	"gamebooks/pkg/storage"
	"gamebooks/pkg/storage/gorm"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"gorm.io/driver/sqlite"
)

type Container struct {
	Repo     repo.Repo
	Storage  storage.Storage
	Executor *executor.Executor
	Markdown goldmark.Markdown
}

func New() Container {
	p := executor.New()
	r := repo.NewWithLiveReload(p)

	gd := sqlite.Open("data.db")
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

	return Container{
		Executor: p,
		Markdown: m,
		Repo:     r,
		Storage:  s,
	}
}
