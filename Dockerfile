FROM golang:1.6.2-alpine

MAINTAINER TFG Co <backend@tfgco.com>

EXPOSE 8080

RUN apk update
RUN apk add git bash

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
ENV MQTTBOT_API_TLS false
ENV MQTTBOT_API_CERTFILE ./misc/example.crt
ENV MQTTBOT_API_KEYFILE ./misc/example.key

CMD ./start_docker.sh
