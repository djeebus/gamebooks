package web

import (
	"gamebooks/pkg/executor"
	bookRepo "gamebooks/pkg/repo"
	"gamebooks/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func New(g bookRepo.Game, storage storage.Storage, player *executor.Player) (*echo.Echo, error) {
	e := echo.New()
	e.Debug = true
	e.HideBanner = true

	e.Use(logRequests)
	e.Use(addSessionID)

	v, err := newViews(g, storage, player)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create views")
	}

	e.GET("/p/:book", v.getBook)
	e.GET("/p/:book/:page", v.getPage)

	e.GET("/", v.index)

	return e, nil
}
