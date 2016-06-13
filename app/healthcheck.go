package app

import (
	"github.com/kataras/iris"
	"github.com/spf13/viper"
)

//HealthCheckHandler is the handler responsible for validating that the app is still up
func HealthCheckHandler(app *App) func(c *iris.Context) {
	return func(c *iris.Context) {
		workingString := viper.GetString("healthcheck.workingText")
		c.SetStatusCode(iris.StatusOK)
		c.Write(workingString)
	}
}
