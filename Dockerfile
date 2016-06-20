FROM golang:1.6.2-alpine

MAINTAINER TFG Co <backend@tfgco.com>

EXPOSE 8080

RUN apk update
RUN apk add git

RUN go get -u github.com/Masterminds/glide/...

ADD . /go/src/github.com/topfreegames/mqttbot

WORKDIR /go/src/github.com/topfreegames/mqttbot
RUN glide install
RUN go install github.com/topfreegames/mqttbot

ENV MQTTBOT_MQTTSERVER_HOST localhost
ENV MQTTBOT_MQTTSERVER_PORT 1883
ENV MQTTBOT_MQTTSERVER_USER admin
ENV MQTTBOT_MQTTSERVER_PASS admin

ENV MQTTBOT_ELASTICSEARCH_HOST http://localhost:9200
ENV MQTTBOT_ELASTICSEARCH_SNIFF false

ENV MQTTBOT_REDIS_HOST localhost
ENV MQTTBOT_REDIS_PORT 6379
ENV MQTTBOT_REDIS_PASSWORD pass

CMD /go/bin/mqttbot start --bind 0.0.0.0 --port 8080 --config /go/src/github.com/topfreegames/mqttbot/config/local.yaml
