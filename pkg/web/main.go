package web

import (
	bookRepo "gamebooks/pkg/repo"
	"gamebooks/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

func New(g bookRepo.Game, storage storage.Storage) (*echo.Echo, error) {
	e := echo.New()
	e.Debug = true
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(addSessionID)

	v, err := newViews(g, storage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create views")
	}

	e.GET("/p/:book", v.getBook)
	e.GET("/p/:book/:page", v.getPage)

	e.GET("/", v.index)

	return e, nil
}
