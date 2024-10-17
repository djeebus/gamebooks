package web

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

const authCookieName = ".auth"
const authCookieDuration = time.Hour * 24 * 30 // one month

func addSessionID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(authCookieName)
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				return errors.Wrap(err, "invalid cookie")
			}

			cookie = &http.Cookie{
				Name:    authCookieName,
				Value:   uuid.NewString(),
				Expires: time.Now().Add(authCookieDuration),
			}
			c.SetCookie(cookie)
		}

		return next(c)
	}
}
