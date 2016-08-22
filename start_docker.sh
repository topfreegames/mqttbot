#!/bin/sh

if [ "$MQTTBOT_API_TLS" == "true" ]
then
  echo -e $MQTTBOT_API_CERT > $MQTTBOT_API_CERTFILE
  echo -e $MQTTBOT_API_KEY > $MQTTBOT_API_KEYFILE
  /go/bin/mqttbot start --bind 0.0.0.0 --port 4443 --config /go/src/github.com/topfreegames/mqttbot/config/local.yaml
else
  /go/bin/mqttbot start --bind 0.0.0.0 --port 5000 --config /go/src/github.com/topfreegames/mqttbot/config/local.yaml
fi
