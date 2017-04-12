#!/bin/sh
COUNT=`curl -XPOST -d'{"query":{"bool":{"must":[{"range":{"timestamp":{"gte":"now-1m"}}}]}},"size":0}' "${MQTTBOT_ELASTICSEARCH_HOST}/chat/_search" | jq .hits.total`
if [[ $COUNT -eq 0 ]]; then
  exit 1
fi
