// mqttbot
// https://github.com/topfreegames/mqttbot
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package app

import (
	"net/http"
	"testing"
	"time"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/topfreegames/mqttbot/es"
	"github.com/topfreegames/mqttbot/redisclient"
)

func TestHistoryHandler(t *testing.T) {
	g := Goblin(t)

	// special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("History Handler", func() {
		g.It("It should return 404 if the user is not authorized into the topic", func() {
			a := GetDefaultTestApp()
			res := Get(a, "/history/chat/teste?userid=test:teste", t)
			g.Assert(res.Raw().StatusCode).Equal(http.StatusForbidden)
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
			res := GetWithQuery(a, "/history/chat/teste", "userid", "test:teste", t)
			g.Assert(res.Raw().StatusCode).Equal(http.StatusOK)
		})
	})
}
