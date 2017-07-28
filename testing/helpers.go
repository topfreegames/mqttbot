// mqttbot
// https://github.com/topfreegames/mqttbot
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/topfreegames/mqttbot/app"
)

// GetDefaultTestApp retrieve a default app for testing purposes
func GetDefaultTestApp() *app.App {
	viper.SetConfigFile("../config/test.yaml")
	app := app.GetApp("0.0.0.0", 8888, true)

	return app
}

// Get implements the GET http verb for testing purposes
func Get(app *app.App, url string, t *testing.T) (int, string) {
	return doRequest(app, "GET", url, "")
}

/*
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
*/

func doRequest(app *app.App, method, url, body string) (int, string) {
	// app.Engine.SetHandler(app.API)
	ts := httptest.NewServer(app.API)
	defer ts.Close()

	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", ts.URL, url), reader)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	client := &http.Client{}
	res, err := client.Do(req)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	if res != nil {
		defer res.Body.Close()
	}

	b, err := ioutil.ReadAll(res.Body)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	return res.StatusCode, string(b)
}
