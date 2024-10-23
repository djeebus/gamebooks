package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"runtime/debug"
)

func logRequests(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			r := recover()
			if r != nil {
				stack := debug.Stack()
				fmt.Println(string(stack))
			}
		}()

		if err := next(c); err != nil {
			return c.String(500, "internal server error")
		}

		return nil
	}
}

func recordErrors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			msg := fmt.Sprintf("unhandled server panic:\n%s", err)
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
				msg := fmt.Sprintf("%v\n\n%s", r, stack)
				log.Error().Msg(msg)
				err = c.String(500, msg)
			}
		}()

		err = next(c)
		return err
	}
}
