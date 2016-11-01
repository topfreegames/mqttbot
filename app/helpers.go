// mqttbot
// https://github.com/topfreegames/mqttbot
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package app

import (
	"github.com/labstack/echo"
	newrelic "github.com/newrelic/go-agent"
)

//GetTX returns new relic transaction
func GetTX(c echo.Context) newrelic.Transaction {
	tx := c.Get("txn")
	if tx == nil {
		return nil
	}

	return tx.(newrelic.Transaction)
}

//WithSegment adds a segment to new relic transaction
func WithSegment(name string, c echo.Context, f func() error) error {
	tx := GetTX(c)
	if tx == nil {
		return f()
	}
	segment := newrelic.StartSegment(tx, name)
	defer segment.End()
	return f()
}
