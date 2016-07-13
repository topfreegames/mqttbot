package app

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/kataras/iris"
	"github.com/topfreegames/mqttbot/es"
	"github.com/topfreegames/mqttbot/logger"
	"gopkg.in/topfreegames/elastic.v2"
)

// PayloadStruct contains the fields of a payload
type PayloadStruct struct {
	From      string `json:"from"`
	Message   string `json:"message"`
	Timestamp int32  `json:"timestamp"`
	Id        string `json:"id"`
}

type Message struct {
	Payload PayloadStruct `json:"payload"`
	Topic   string        `json:"topic"`
}

// HistoryHandler is the handler responsible for sending the rooms history to the player
func HistoryHandler(app *App) func(c *iris.Context) {
	return func(c *iris.Context) {
		esclient := es.GetESClient()
		topic := c.Param("topic")[1:len(c.Param("topic"))]
		userId := c.URLParam("userid")
		from, _ := c.URLParamInt("from")
		limit, _ := c.URLParamInt("limit")
		if limit == 0 {
			limit = 10
		}
		rc := app.RedisClient.Pool.Get()
		rc.Send("MULTI")
		rc.Send("GET", userId)
		rc.Send("GET", fmt.Sprintf("%s-%s", userId, topic))
		r, err := rc.Do("EXEC")
		if err != nil {
			logger.Logger.Error(err.Error())
			c.SetStatusCode(iris.StatusInternalServerError)
			return
		}
		redisResults := (r.([]interface{}))
		if redisResults[0] != nil && redisResults[1] != nil {
			termQuery := elastic.NewQueryStringQuery(fmt.Sprintf("topic:\"%s\"", topic))
			searchResults, err := esclient.Search().Index("chat").Query(termQuery).
				Sort("payload.timestamp", false).From(from).Size(limit).Do()
			if err != nil {
				logger.Logger.Error(err.Error())
				c.SetStatusCode(iris.StatusInternalServerError)
				return
			}
			payloads := []PayloadStruct{}
			var ttyp Message
			for _, item := range searchResults.Each(reflect.TypeOf(ttyp)) {
				if t, ok := item.(Message); ok {
					payloads = append(payloads, t.Payload)
				}
			}
			jsonPayloads, _ := json.Marshal(payloads)
			c.Write(fmt.Sprintf("%s", jsonPayloads))
			c.SetStatusCode(iris.StatusOK)
		} else {
			c.SetStatusCode(iris.StatusForbidden)
			return
		}
	}
}
