package web

import (
	"gamebooks/pkg/executor"
	bookRepo "gamebooks/pkg/repo"
	"gamebooks/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"
)

func New(g bookRepo.Game, storage storage.Storage, player *executor.Player, log zerolog.Logger) (*echo.Echo, error) {
	e := echo.New()
	e.Debug = true
	e.HideBanner = true
	e.Logger = lecho.From(log)

	e.Use(middleware.Logger())
	e.Use(recordPanics)
	e.Use(recordErrors)
	e.Use(addSessionID)

	v, err := newViews(g, storage, player)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create views")
	}

	e.GET("/game", v.gameView)
	e.POST("/game", v.gameNext)

	e.GET("/", v.listBooks)
	e.GET("/start/:bookID", v.selectBook)

	return e, nil
}
