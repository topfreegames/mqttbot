package app

import "testing"

func GetAppTest(t *testing.T) {
	app := GetApp("127.0.0.1", 9999, false)
	if app.Port != 9999 || app.Host != "127.0.0.1" {
		t.Fail()
	}
}
