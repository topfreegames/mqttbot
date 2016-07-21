// mqttbot
// https://github.com/topfreegames/mqttbot
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package app

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/spf13/viper"
)

func GetDefaultTestApp() *App {
	viper.SetConfigFile("../config/test.yaml")
	app := GetApp("0.0.0.0", 8888, true)
	return app
}

func Get(app *App, url string, t *testing.T) *httpexpect.Response {
	req := sendRequest(app, "GET", url, t)
	return req.Expect()
}

func GetWithQuery(app *App, url string, queryKey string, queryValue string, t *testing.T) *httpexpect.Response {

	srv := app.Api.Servers.Main()

	if srv == nil { // maybe the user called this after .Listen/ListenTLS/ListenUNIX, the t
		srv = app.Api.ListenVirtual(app.Api.Config.Tester.ListeningAddr)
	}

	handler := srv.Handler
	e := httpexpect.WithConfig(httpexpect.Config{
		Reporter: httpexpect.NewAssertReporter(t),
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(handler),
		},
	})

	return e.GET(url).WithQuery(queryKey, queryValue).Expect()
}

func sendRequest(app *App, method, url string, t *testing.T) *httpexpect.Request {

	srv := app.Api.Servers.Main()

	if srv == nil { // maybe the user called this after .Listen/ListenTLS/ListenUNIX, the t
		srv = app.Api.ListenVirtual(app.Api.Config.Tester.ListeningAddr)
	}

	handler := srv.Handler

	e := httpexpect.WithConfig(httpexpect.Config{
		Reporter: httpexpect.NewAssertReporter(t),
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(handler),
		},
	})

	return e.Request(method, url)
}
