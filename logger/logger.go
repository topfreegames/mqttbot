package logger

import (
	"os"

	"github.com/op/go-logging"
)

var Logger = logging.MustGetLogger("ggas")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func SetupLogger(logLevel string) {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	level, err := logging.LogLevel(logLevel)
	if err != nil {
		panic(err)
	}
	backendLeveled.SetLevel(level, "")
	logging.SetBackend(backendLeveled)
}
