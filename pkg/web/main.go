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

func New(g bookRepo.Game, storage storage.Storage, executor *executor.Executor, log zerolog.Logger) (*echo.Echo, error) {
	e := echo.New()
	e.Debug = true
	e.HideBanner = true
	e.Logger = lecho.From(log)

	e.Use(middleware.Logger())
	e.Use(recordPanics)
	e.Use(recordErrors)
	e.Use(addSessionID)

	v, err := newViews(g, storage, executor)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create views")
	}

	e.GET("/", v.gameView)
	e.POST("/", v.gameNext)

	return e, nil
}
