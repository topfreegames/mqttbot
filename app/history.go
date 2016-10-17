package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/topfreegames/mqttbot/es"
	"github.com/topfreegames/mqttbot/logger"
	"gopkg.in/olivere/elastic.v3"
)

// Message represents a chat message
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
		userID := c.QueryParam("userid")
		from, err := strconv.Atoi(c.QueryParam("from"))
		limit, err := strconv.Atoi(c.QueryParam("limit"))

		if limit == 0 {
			limit = 10
		}

		logger.Logger.Debugf("user %s is asking for history for topic %s with args from=%d and limit=%d", userID, topic, from, limit)
		rc := app.RedisClient.Pool.Get()
		defer rc.Close()
		rc.Send("MULTI")
		rc.Send("GET", userID)
		rc.Send("GET", fmt.Sprintf("%s-%s", userID, topic))
		r, err := rc.Do("EXEC")
		if err != nil {
			return err
		}
		redisResults := (r.([]interface{}))
		if redisResults[0] != nil && redisResults[1] != nil {
			boolQuery := elastic.NewBoolQuery()

			termQuery := elastic.NewTermQuery("topic", topic)
			boolQuery.Must(termQuery)

			searchResults, err := esclient.Search().Index("chat").Query(boolQuery).
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
			resStr, err := json.Marshal(messages)
			if err != nil {
				return err
			}

			logger.Logger.Debugf(
				"responded to user %s history for topic %s with args from=%d and limit=%d with code=%s and message=%s",
				userID, topic, from, limit, http.StatusOK, echo.ErrUnauthorized.Message, string(resStr),
			)
			return c.JSON(http.StatusOK, messages)
		}
		logger.Logger.Debugf(
			"responded to user %s history for topic %s with args from=%d and limit=%d with code=%s and message=%s",
			userID, topic, from, limit, echo.ErrUnauthorized.Code, echo.ErrUnauthorized.Message,
		)
		return c.String(echo.ErrUnauthorized.Code, echo.ErrUnauthorized.Message)
	}
}

// HistorySinceHandler is the handler responsible for sending the rooms history to the player based in a initial date
func HistorySinceHandler(app *App) func(c echo.Context) error {
	return func(c echo.Context) error {
		esclient := es.GetESClient()
		c.Set("route", "HistorySince")
		topic := c.ParamValues()[0]
		userID := c.QueryParam("userid")
		from, err := strconv.Atoi(c.QueryParam("from"))
		limit, err := strconv.Atoi(c.QueryParam("limit"))
		since := c.QueryParam("since")

		if limit == 0 {
			limit = 10
		}

		logger.Logger.Debugf("user %s is asking for history for topic %s with args from=%d, limit=%d and since=%s", userID, topic, from, limit, since)
		rc := app.RedisClient.Pool.Get()
		defer rc.Close()
		rc.Send("MULTI")
		rc.Send("GET", userID)
		rc.Send("GET", fmt.Sprintf("%s-%s", userID, topic))
		r, err := rc.Do("EXEC")
		if err != nil {
			return err
		}

		redisResults := (r.([]interface{}))
		if redisResults[0] != nil && redisResults[1] != nil {
			boolQuery := elastic.NewBoolQuery()

			termQuery := elastic.NewTermQuery("topic", topic)
			boolQuery.Must(termQuery)

			if since != "" {
				rangeQuery := elastic.NewRangeQuery("timestamp").
					From(since).
					To(nil).
					IncludeLower(true).
					IncludeUpper(true)

				boolQuery.Must(rangeQuery, termQuery)
			} else {
				boolQuery.Must(termQuery)
			}

			searchResults, err := esclient.Search().Index("chat").Query(boolQuery).
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
		}
		return c.String(echo.ErrUnauthorized.Code, echo.ErrUnauthorized.Message)
	}
}
