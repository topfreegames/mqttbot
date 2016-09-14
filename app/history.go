package app

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"gopkg.in/olivere/elastic.v3"

	"github.com/labstack/echo"
	"github.com/topfreegames/mqttbot/es"
	"github.com/topfreegames/mqttbot/logger"
)

type Message struct {
	Timestamp time.Time `json:"timestamp"`
	Payload   string    `json:"payload"`
	Topic     string    `json:"topic"`
}

// HistoryHandler is the handler responsible for sending the rooms history to the player
func HistoryHandler(app *App) func(c echo.Context) error {
	return func(c echo.Context) error {
		esclient := es.GetESClient()
		c.Set("route", "History")
		topic := c.ParamValues()[0]
		userId := c.QueryParam("userid")
		from, err := strconv.Atoi(c.QueryParam("from"))
		limit, err := strconv.Atoi(c.QueryParam("limit"))

		if limit == 0 {
			limit = 10
		}

		logger.Logger.Debugf("user %s is asking for history for topic %s with args from=%d and limit=%d", userId, topic, from, limit)
		rc := app.RedisClient.Pool.Get()
		defer rc.Close()
		rc.Send("MULTI")
		rc.Send("GET", userId)
		rc.Send("GET", fmt.Sprintf("%s-%s", userId, topic))
		r, err := rc.Do("EXEC")
		if err != nil {
			return err
		}
		redisResults := (r.([]interface{}))
		if redisResults[0] != nil && redisResults[1] != nil {
			termQuery := elastic.NewQueryStringQuery(fmt.Sprintf("topic:\"%s\"", topic))
			searchResults, err := esclient.Search().Index("chat").Query(termQuery).
				Sort("timestamp", false).From(from).Size(limit).Do()
			if err != nil {
				return err
			}
			messages := []Message{}
			var ttyp Message
			for _, item := range searchResults.Each(reflect.TypeOf(ttyp)) {
				if t, ok := item.(Message); ok {
					messages = append(messages, t)
				}
			}
			return c.JSON(http.StatusOK, messages)
		} else {
			return c.String(echo.ErrUnauthorized.Code, echo.ErrUnauthorized.Message)
		}
	}
}
