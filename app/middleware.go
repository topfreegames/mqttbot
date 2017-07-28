package app

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

//VersionMiddleware automatically adds a version header to response
func VersionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, fmt.Sprintf("mqttbot/v%s", VERSION))
		return next(c)
	}
}

//NewRecoveryMiddleware returns a configured middleware
func NewRecoveryMiddleware(onError func(interface{}, []byte)) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		OnError: onError,
	}
}

//RecoveryMiddleware recovers from errors in Echo
type RecoveryMiddleware struct {
	OnError func(interface{}, []byte)
}

//Serve executes on error handler when errors happen
func (r RecoveryMiddleware) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				if r.OnError != nil {
					r.OnError(err, debug.Stack())
				}

				if eError, ok := err.(error); ok {
					c.Error(eError)
				} else {
					eError = fmt.Errorf(fmt.Sprintf("%v", err))
					c.Error(eError)
				}
			}
		}()
		return next(c)
	}
}

//LoggerMiddleware is responsible for logging to Zap all requests
type LoggerMiddleware struct {
	Logger log.FieldLogger
}

// Serve serves the middleware
func (l LoggerMiddleware) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lg := l.Logger.WithField("source", "request")

		//all except latency to string
		var ip, method, path string
		var status int
		var latency time.Duration
		var startTime, endTime time.Time

		path = c.Path()
		method = c.Request().Method

		startTime = time.Now()

		result := next(c)

		//no time.Since in order to format it well after
		endTime = time.Now()
		latency = endTime.Sub(startTime)

		status = c.Response().Status
		ip = c.Request().RemoteAddr

		route := c.Get("route")
		if route == nil {
			lg.Debug("Route does not have route set in ctx")
			return result
		}

		reqLog := lg.WithFields(log.Fields{
			"route":      route.(string),
			"endTime":    endTime,
			"statusCode": status,
			"latency":    latency,
			"ip":         ip,
			"method":     method,
			"path":       path,
		})

		//request failed
		if status > 399 && status < 500 {
			reqLog.Warn("Request failed.")
			return result
		}

		//request is ok, but server failed
		if status > 499 {
			reqLog.Error("Response failed.")
			return result
		}
		//Everything went ok
		reqLog.Info("Request successful.")
		return result
	}
}

// NewLoggerMiddleware returns the logger middleware
func NewLoggerMiddleware(theLogger log.FieldLogger) *LoggerMiddleware {
	l := &LoggerMiddleware{Logger: theLogger}
	return l
}

//SentryMiddleware is responsible for sending all exceptions to sentry
type SentryMiddleware struct {
	App *App
}

// Serve serves the middleware
func (s SentryMiddleware) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			tags := map[string]string{
				"source": "app",
				"type":   "Internal server error",
				"url":    c.Request().URL.String(),
				"status": fmt.Sprintf("%d", c.Response().Status),
			}
			raven.CaptureError(err, tags)
		}
		return err
	}
}

//NewSentryMiddleware returns a new sentry middleware
func NewSentryMiddleware(app *App) *SentryMiddleware {
	return &SentryMiddleware{
		App: app,
	}
}

//NewNewRelicMiddleware returns the logger middleware
func NewNewRelicMiddleware(app *App, theLogger log.FieldLogger) *NewRelicMiddleware {
	l := &NewRelicMiddleware{App: app, Logger: theLogger}
	return l
}

//NewRelicMiddleware is responsible for logging to Zap all requests
type NewRelicMiddleware struct {
	App    *App
	Logger log.FieldLogger
}

// Serve serves the middleware
func (nr *NewRelicMiddleware) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		route := c.Path()
		txn := nr.App.NewRelic.StartTransaction(route, nil, nil)
		c.Set("txn", txn)
		defer func() {
			c.Set("txn", nil)
			txn.End()
		}()

		err := next(c)
		if err != nil {
			txn.NoticeError(err)
			return err
		}

		return nil
	}
}
