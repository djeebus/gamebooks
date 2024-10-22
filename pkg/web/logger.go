package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
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

		return next(c)
	}
}
