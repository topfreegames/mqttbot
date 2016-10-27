package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

		defaultLimit := 10
		if limitFromEnv := os.Getenv("HISTORY_LIMIT"); limitFromEnv != "" {
			defaultLimit, err = strconv.Atoi(limitFromEnv)
		}
		if limit == 0 {
			limit = defaultLimit
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
				"responded to user %s history for topic %s with args from=%d and limit=%d with code=%d and message=%s",
				userID, topic, from, limit, http.StatusOK, string(resStr),
			)
			return c.JSON(http.StatusOK, messages)
		}
		logger.Logger.Debugf(
			"responded to user %s history for topic %s with args from=%d and limit=%d with code=%d and message=%s",
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
		since, err := strconv.ParseInt(c.QueryParam("since"), 10, 64)

		defaultLimit := 10
		if limitFromEnv := os.Getenv("HISTORYSINCE_LIMIT"); limitFromEnv != "" {
			defaultLimit, err = strconv.Atoi(limitFromEnv)
		}
		if limit == 0 {
			limit = defaultLimit
		}

		logger.Logger.Debugf("user %s is asking for history for topic %s with args from=%d, limit=%d and since=%d", userID, topic, from, limit, since)
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
			rangeQuery := elastic.NewRangeQuery("timestamp").
				From(since * 1000). // FIXME: The client should send time in milliseconds
				To(nil).
				IncludeLower(true).
				IncludeUpper(true)
			boolQuery.Must(rangeQuery, termQuery)

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
				"responded to user %s history for topic %s with args from=%d limit=%d and since=%d with code=%d and message=%s",
				userID, topic, from, limit, since, http.StatusOK, string(resStr),
			)
			return c.JSON(http.StatusOK, messages)
		}
		logger.Logger.Debugf(
			"responded to user %s history for topic %s with args from=%d limit=%d and since=%d with code=%d and message=%s",
			userID, topic, from, limit, since, echo.ErrUnauthorized.Code, echo.ErrUnauthorized.Message,
		)
		return c.String(echo.ErrUnauthorized.Code, echo.ErrUnauthorized.Message)
	}
}
