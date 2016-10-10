package app

import "github.com/labstack/echo"

// NotFoundHandler is the handler responsible for responding when no resource was found
func NotFoundHandler(app *App) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.String(echo.ErrNotFound.Code, echo.ErrNotFound.Message)
	}
}
