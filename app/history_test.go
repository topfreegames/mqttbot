// mqttbot
// https://github.com/topfreegames/mqttbot
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
	. "github.com/topfreegames/mqttbot/app"
	"github.com/topfreegames/mqttbot/es"
	"github.com/topfreegames/mqttbot/redisclient"
	. "github.com/topfreegames/mqttbot/testing"
)

func refreshIndex() {
	_, err := http.Post("http://localhost:9123/_refresh", "application/json", bytes.NewBufferString("{}"))
	Expect(err).To(BeNil())
}

func msToTime(ms int64) time.Time {
	return time.Unix(0, ms*int64(time.Millisecond))
}

func TestHistoryHandler(t *testing.T) {
	g := Goblin(t)

	// special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("History", func() {
		esclient := es.GetESClient()
		g.BeforeEach(func() {
			esclient.DeleteIndex("chat")
			refreshIndex()
			// esclient.CreateIndex("chat")
		})

		g.Describe("History Handler", func() {
			g.It("It should return 401 if the user is not authorized into the topic", func() {
				a := GetDefaultTestApp()
				testId := strings.Replace(uuid.NewV4().String(), "-", "", -1)
				path := fmt.Sprintf("/history/chat/test_%s?userid=test:test", testId)
				status, _ := Get(a, path, t)
				g.Assert(status).Equal(http.StatusUnauthorized)
			})

			g.It("It should return 200 if the user is authorized into the topic", func() {
				a := GetDefaultTestApp()
				testId := strings.Replace(uuid.NewV4().String(), "-", "", -1)
				topic := fmt.Sprintf("chat/test_%s", testId)
				authStr := fmt.Sprintf("test:test-%s", topic)
				rc := redisclient.GetRedisClient("localhost", 4444, "")
				_, err := rc.Pool.Get().Do("set", "test:test", "lalala")
				_, err = rc.Pool.Get().Do("set", authStr, 2)
				Expect(err).To(BeNil())

				testMessage := Message{
					Timestamp: time.Now(),
					Payload:   "{\"test1\":\"test2\"}",
					Topic:     topic,
				}
				_, err = esclient.Index().Index("chat").Type("message").BodyJson(testMessage).Do()
				Expect(err).To(BeNil())

				refreshIndex()
				path := fmt.Sprintf("/history/%s?userid=test:test", topic)
				status, body := Get(a, path, t)
				g.Assert(status).Equal(http.StatusOK)

				var messages []Message
				err = json.Unmarshal([]byte(body), &messages)
				Expect(err).To(BeNil())
			})
		})

		g.Describe("History Since Handler", func() {
			g.It("It should return 401 if the user is not authorized into the topic", func() {
				a := GetDefaultTestApp()
				testId := strings.Replace(uuid.NewV4().String(), "-", "", -1)
				path := fmt.Sprintf("/history/chat/test_%s?userid=test:test", testId)
				status, _ := Get(a, path, t)
				g.Assert(status).Equal(http.StatusUnauthorized)
			})

			g.It("It should return 200 if the user is authorized into the topic", func() {
				a := GetDefaultTestApp()
				testId := strings.Replace(uuid.NewV4().String(), "-", "", -1)
				topic := fmt.Sprintf("chat/test_%s", testId)
				authStr := fmt.Sprintf("test:test-%s", topic)

				rc := redisclient.GetRedisClient("localhost", 4444, "")
				_, err := rc.Pool.Get().Do("set", "test:test", "lalala")
				_, err = rc.Pool.Get().Do("set", authStr, 2)
				Expect(err).To(BeNil())

				testMessage := Message{
					Timestamp: time.Now(),
					Payload:   "{\"test1\":\"test2\"}",
					Topic:     topic,
				}

				_, err = esclient.Index().Index("chat").Type("message").BodyJson(testMessage).Do()
				Expect(err).To(BeNil())

				refreshIndex()

				path := fmt.Sprintf("/history/%s?userid=test:test", topic)
				status, body := Get(a, path, t)
				g.Assert(status).Equal(http.StatusOK)

				var messages []Message
				err = json.Unmarshal([]byte(body), &messages)
				Expect(err).To(BeNil())
			})

			g.It("It should return 200 if the user is authorized into the topic", func() {
				a := GetDefaultTestApp()
				testId := strings.Replace(uuid.NewV4().String(), "-", "", -1)
				topic := fmt.Sprintf("chat/test_%s", testId)
				authStr := fmt.Sprintf("test:test-%s", topic)

				rc := redisclient.GetRedisClient("localhost", 4444, "")
				_, err := rc.Pool.Get().Do("set", "test:test", "lalala")
				_, err = rc.Pool.Get().Do("set", authStr, 2)
				Expect(err).To(BeNil())

				testMessage := Message{
					Timestamp: time.Now(),
					Payload:   "{\"test1\":\"test2\"}",
					Topic:     topic,
				}

				path := fmt.Sprintf(
					"/historysince/%s?userid=test:test&since=%d",
					topic, (time.Now().UnixNano() / 1000000), // now
				)
				_, err = esclient.Index().Index("chat").Type("message").BodyJson(testMessage).Do()
				Expect(err).To(BeNil())

				// Update indexes
				refreshIndex()

				status, body := Get(a, path, t)
				g.Assert(status).Equal(http.StatusOK)

				var messages []Message
				err = json.Unmarshal([]byte(body), &messages)
				Expect(err).To(BeNil())
				Expect(len(messages)).To(Equal(1))
				var message Message
				for i := 0; i < len(messages); i++ {
					message = messages[i]
					Expect(message.Topic).To(Equal(topic))
				}
			})

			g.It("It should return 200 if the user is authorized into the topic", func() {
				a := GetDefaultTestApp()
				testId := strings.Replace(uuid.NewV4().String(), "-", "", -1)
				topic := fmt.Sprintf("chat/test_%s", testId)
				authStr := fmt.Sprintf("test:test-%s", topic)
				rc := redisclient.GetRedisClient("localhost", 4444, "")
				_, err := rc.Pool.Get().Do("set", "test:test", "lalala")
				_, err = rc.Pool.Get().Do("set", authStr, 2)
				Expect(err).To(BeNil())

				now := time.Now().UnixNano() / 1000000
				testMessage := Message{}
				second := int64(1000)
				baseTime := now - (second * 70)
				for i := 0; i <= 30; i++ {
					messageTime := baseTime + 1*second
					testMessage = Message{
						Timestamp: msToTime(messageTime),
						Payload:   "{\"test1\":\"test2\"}",
						Topic:     topic,
					}
					_, err = esclient.Index().Index("chat").Type("message").BodyJson(testMessage).Do()
					Expect(err).To(BeNil())
				}

				// Update indexes
				refreshIndex()

				path := fmt.Sprintf(
					"/historysince/%s?userid=test:test&since=%d&limit=%d&from=%d",
					topic, baseTime, 10, 0,
				)

				status, body := Get(a, path, t)
				g.Assert(status).Equal(http.StatusOK)

				var messages []Message
				err = json.Unmarshal([]byte(body), &messages)
				Expect(err).To(BeNil())
				Expect(len(messages)).To(Equal(10))
				var message Message
				for i := 0; i < len(messages); i++ {
					message = messages[i]
					Expect(message.Topic).To(Equal(topic))
				}
			})
		})

		g.It("It should return 200 if the user is authorized into the topic", func() {
			a := GetDefaultTestApp()
			testId := strings.Replace(uuid.NewV4().String(), "-", "", -1)
			topic := fmt.Sprintf("chat/test_%s", testId)
			authStr := fmt.Sprintf("test:test-%s", topic)
			rc := redisclient.GetRedisClient("localhost", 4444, "")
			_, err := rc.Pool.Get().Do("set", "test:test", "lalala")
			_, err = rc.Pool.Get().Do("set", authStr, 2)
			Expect(err).To(BeNil())

			startTime := time.Now().UnixNano() / 1000000
			testMessage := Message{}
			for i := 0; i < 3; i++ {
				messageTime := time.Now().UnixNano() / 1000000
				testMessage = Message{
					Timestamp: msToTime(messageTime),
					Payload:   "{\"test1\":\"test2\"}",
					Topic:     topic,
				}
				_, err = esclient.Index().Index("chat").Type("message").BodyJson(testMessage).Do()
				Expect(err).To(BeNil())
			}

			// Sorry bout this =/
			time.Sleep(200 * time.Millisecond)

			// Update indexes
			refreshIndex()

			path := fmt.Sprintf(
				"/historysince/%s?userid=test:test&since=%d&limit=%d&from=%d",
				topic, startTime, 10, 0,
			)

			status, body := Get(a, path, t)
			g.Assert(status).Equal(http.StatusOK)

			var messages []Message
			err = json.Unmarshal([]byte(body), &messages)
			Expect(err).To(BeNil())
			Expect(len(messages)).To(Equal(3))
			var message Message
			for i := 0; i < len(messages); i++ {
				message = messages[i]
				Expect(message.Topic).To(Equal(topic))
			}
		})

		g.It("It should return 200 if the user is authorized into the topic", func() {
			a := GetDefaultTestApp()
			testId := strings.Replace(uuid.NewV4().String(), "-", "", -1)
			topic := fmt.Sprintf("chat/test_%s", testId)
			authStr := fmt.Sprintf("test:test-%s", topic)
			rc := redisclient.GetRedisClient("localhost", 4444, "")
			_, err := rc.Pool.Get().Do("set", "test:test", "lalala")
			_, err = rc.Pool.Get().Do("set", authStr, 2)
			Expect(err).To(BeNil())

			startTime := time.Now().UnixNano() / 1000000
			testMessage := Message{}
			for i := 0; i < 3; i++ {
				messageTime := time.Now().UnixNano() / 1000000
				testMessage = Message{
					Timestamp: msToTime(messageTime),
					Payload:   "{\"test1\":\"test2\"}",
					Topic:     topic,
				}
				_, err = esclient.Index().Index("chat").Type("message").BodyJson(testMessage).Do()
				Expect(err).To(BeNil())
			}

			// Sorry bout this =/
			time.Sleep(200 * time.Millisecond)

			// Update indexes
			refreshIndex()

			path := fmt.Sprintf(
				"/historysince/%s?userid=test:test&since=%d&limit=%d&from=%d",
				topic, startTime, 1, 0,
			)

			status, body := Get(a, path, t)
			g.Assert(status).Equal(http.StatusOK)

			var messages []Message
			err = json.Unmarshal([]byte(body), &messages)
			Expect(err).To(BeNil())
			Expect(len(messages)).To(Equal(1))

			var message Message
			for i := 0; i < len(messages); i++ {
				message = messages[i]
				Expect(message.Topic).To(Equal(topic))
			}
		})
	})
}
