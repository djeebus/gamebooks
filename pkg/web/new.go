package web

import (
	"gamebooks/pkg/executor"
	bookRepo "gamebooks/pkg/repo"
	"gamebooks/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/yuin/goldmark"
	"github.com/ziflex/lecho/v3"
	"net/http"
)

func New(
	repo bookRepo.Repo,
	storage storage.Storage,
	executor *executor.Executor,
	log zerolog.Logger,
	markdown goldmark.Markdown,
) (*echo.Echo, error) {
	e := echo.New()
	e.Debug = true
	e.HideBanner = true
	e.Logger = lecho.From(log)

	e.Use(middleware.Logger())
	e.Use(recordPanics)
	e.Use(recordErrors)
	e.Use(addSessionID)

	v, err := newViews(repo, storage, executor, markdown)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create views")
	}

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/b/")
	})
	e.GET("/b/", v.listBooks)
	e.GET("/b/:bookID", v.gameView)
	e.GET("/b/:bookID/-/clear", v.gameClear)
	e.GET("b/:bookID/-/page/:pageID", v.setPageID)

	return e, nil
}
