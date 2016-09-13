package app

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

//HealthCheckHandler is the handler responsible for validating that the app is still up
func HealthCheckHandler(app *App) func(c echo.Context) error {
	return func(c echo.Context) error {
		c.Set("route", "Healthcheck")
		workingString := viper.GetString("healthcheck.workingText")
		return c.String(http.StatusOK, workingString)
	}
}
