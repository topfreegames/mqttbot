package modules

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"github.com/topfreegames/mqttbot/es"
	"github.com/topfreegames/mqttbot/logger"
	"github.com/yuin/gopher-lua"
	"gopkg.in/olivere/elastic.v3"
)

type PayloadStruct struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

type Message struct {
	Id      string        `json:"id"`
	Payload PayloadStruct `json:"payload"`
	Topic   string        `json:"topic"`
	Date    time.Time     `json:"date"`
}

var ESClient *elastic.Client

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
	ESClient = es.GetESClient()
}

func IndexMessage(L *lua.LState) int {
	topic := L.Get(-2)
	payload := L.Get(-1)
	L.Pop(2)
	var message Message
	json.Unmarshal([]byte(payload.String()), &message)
	message.Topic = topic.String()
	message.Date = time.Now()
	message.Id = uuid.NewV4().String()
	if _, err := ESClient.Index().Index("chat").Type("message").BodyJson(message).Do(); err != nil {
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		L.Push(L.ToNumber(1))
		return 2
	}
	logger.Logger.Debug(fmt.Sprintf("Message persisted: %s", message))
	L.Push(nil)
	L.Push(L.ToNumber(0))
	return 2
}

func QueryMessages(L *lua.LState) int {
	topic := L.Get(-3)
	limit := L.Get(-2)
	start := L.Get(-1)
	L.Pop(3)
	logger.Logger.Debug(fmt.Sprintf("Getting %d messages from topic %s, starting with %d", int(lua.LVAsNumber(limit)), topic.String(), int(lua.LVAsNumber(start))))
	topicSlice := strings.Split(topic.String(), "/")
	topicId := topicSlice[len(topicSlice)-1]
	termQuery := elastic.NewMatchQuery("topic", topicId)
	searchResults, err := ESClient.Search().Index("chat").Query(termQuery).Sort("date", false).From(int(lua.LVAsNumber(start))).Size(int(lua.LVAsNumber(limit))).Do()
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		L.Push(L.ToNumber(1))
		return 2
	}
	var ttyp Message
	for _, item := range searchResults.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(Message); ok {
			logger.Logger.Debug(fmt.Sprintf("Message: %s\n", t))
		}
	}
	return 0
}
