// mqttbot
// https://github.com/topfreegames/mqttbot
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package app_test

import (
	"net/http"
	"testing"
	"time"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
	. "github.com/topfreegames/mqttbot/app"
	"github.com/topfreegames/mqttbot/es"
	"github.com/topfreegames/mqttbot/redisclient"
	. "github.com/topfreegames/mqttbot/testing"
)

func TestHistoryHandler(t *testing.T) {
	g := Goblin(t)

	// special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("History Handler", func() {
		g.It("It should return 401 if the user is not authorized into the topic", func() {
			a := GetDefaultTestApp()
			status, _ := Get(a, "/history/chat/teste?userid=test:teste2", t)
			g.Assert(status).Equal(http.StatusUnauthorized)
		})

		g.It("It should return 200 if the user is authorized into the topic", func() {
			a := GetDefaultTestApp()
			rc := redisclient.GetRedisClient("localhost", 4444, "")
			_, err := rc.Pool.Get().Do("set", "test:teste", "lalala")
			_, err = rc.Pool.Get().Do("set", "test:teste-chat/teste", 2)
			Expect(err).To(BeNil())

			esclient := es.GetESClient()
			testMessage := Message{
				Timestamp: time.Now(),
				Payload:   "{\"test1\":\"test2\"}",
				Topic:     "chat/teste",
			}
			_, err = esclient.Index().Index("chat").Type("message").BodyJson(testMessage).Do()
			Expect(err).To(BeNil())
			status, _ := Get(a, "/history/chat/teste?userid=test:teste", t)
			g.Assert(status).Equal(http.StatusOK)
		})
	})
}
