package modules

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"gopkg.in/olivere/elastic.v5"

	"github.com/layeh/gopher-luar"
	"github.com/topfreegames/mqttbot/es"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/yuin/gopher-lua"
)

type Message struct {
	Timestamp time.Time `json:"timestamp"`
	Payload   string    `json:"payload"`
	Topic     string    `json:"topic"`
}

var esclient *elastic.Client

func PersistenceModuleLoader(L *lua.LState) int {
	configurePersistenceModule()
	mod := L.SetFuncs(L.NewTable(), persistenceModuleExports)
	L.Push(mod)
	return 1
}

var persistenceModuleExports = map[string]lua.LGFunction{
	"index_message":  IndexMessage,
	"query_messages": QueryMessages,
}

func configurePersistenceModule() {
	esclient = es.GetESClient()
}

func IndexMessage(L *lua.LState) int {
	topic := L.Get(-2)
	payload := L.Get(-1)
	L.Pop(2)
	message := Message{}
	message.Timestamp = time.Now()
	message.Payload = payload.String()
	message.Topic = topic.String()
	if _, err := esclient.Index().Index("chat").Type("message").BodyJson(message).Do(context.TODO()); err != nil {
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		L.Push(L.ToNumber(1))
		return 2
	}
	logger.Logger.Debug(fmt.Sprintf("Message persisted: %v", message))
	L.Push(lua.LNil)
	L.Push(L.ToNumber(0))
	return 2
}

func QueryMessages(L *lua.LState) int {
	topic := L.Get(-3)
	limit := L.Get(-2)
	start := L.Get(-1)
	L.Pop(3)
	logger.Logger.Debug(fmt.Sprintf("Getting %d messages from topic %s, starting with %d",
		int(lua.LVAsNumber(limit)), topic.String(), int(lua.LVAsNumber(start))))
	termQuery := elastic.NewQueryStringQuery(fmt.Sprintf("topic:\"%s\"", topic.String()))
	searchResults, err := esclient.Search().Index("chat").Query(termQuery).
		Sort("timestamp", false).From(int(lua.LVAsNumber(start))).
		Size(int(lua.LVAsNumber(limit))).Do(context.TODO())
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		L.Push(L.ToNumber(1))
		return 2
	}
	messages := []Message{}
	var ttyp Message
	for _, item := range searchResults.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(Message); ok {
			logger.Logger.Debug(fmt.Sprintf("Message: %s\n", t))
			messages = append(messages, t)
		}
	}
	L.Push(lua.LNil)
	L.Push(luar.New(L, messages))
	return 2
}
