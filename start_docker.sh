#!/bin/sh

if [ "$MQTTBOT_API_TLS" == "true" ]
then
  echo -e $MQTTBOT_API_CERT > $MQTTBOT_API_CERTFILE
  echo -e $MQTTBOT_API_KEY > $MQTTBOT_API_KEYFILE
  PORT=4443
else
  PORT=5000
fi

/go/bin/mqttbot start --bind 0.0.0.0 --port ${PORT} --config ${MQTTBOT_CONFIG_FILE}
