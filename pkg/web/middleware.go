package web

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"reflect"
	"runtime/debug"
	"time"
)

const authCookieName = ".auth"
const authCookieDuration = time.Hour * 24 * 30 // one month

func addSessionID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := c.Cookie(authCookieName)
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				return errors.Wrap(err, "invalid cookie")
			}

			cookie := &http.Cookie{
				Name:    authCookieName,
				Value:   uuid.NewString(),
				Expires: time.Now().Add(authCookieDuration),
			}
			c.SetCookie(cookie)
		}

		return next(c)
	}
}

func getUserID(c echo.Context) string {
	result, err := c.Cookie(authCookieName)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve cookie")
		return ""
	}

	return result.Value
}

func recordErrors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			log.Error().Err(err).Str("err_type", reflect.TypeOf(err).String()).Stack().Msg("internal server error")
			msg := fmt.Sprintf("internal server error:\n%s", err)
			return c.String(500, msg)
		}

		return nil
	}
}

func recordPanics(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				stack := debug.Stack()
				msg := fmt.Sprintf("unhandled server panic: %v\n\n%s", r, stack)
				log.Error().Msg(msg)
				err = c.String(500, msg)
			}
		}()

		err = next(c)
		return err
	}
}
