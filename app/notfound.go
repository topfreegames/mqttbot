package app

import "github.com/labstack/echo"

//HealthCheckHandler is the handler responsible for validating that the app is still up
func NotFoundHandler(app *App) func(c echo.Context) error {
	return func(c echo.Context) error {
		return echo.ErrNotFound
	}
}
