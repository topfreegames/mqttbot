package models

import (
	"gopkg.in/olivere/elastic.v3"
	"strconv"
	"time"
)

type Message struct {
	Message string
	Topic   string
}

func (m *Message) Index(client *elastic.Client) error {
	_, err := client.Index().
		Index("chat").
		Type("message").
		Timestamp(strconv.FormatInt(time.Now().Unix()*1000, 10)).
		TTL("2d").
		BodyJson(m).
		Do()
	return err
}

func GetMessages(fromTimestamp int, toTimestamp int, limit int) *[]Message {
	return nil
}
