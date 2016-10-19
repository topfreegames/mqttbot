#!/bin/sh

echo 'Create chat index template: '
curl -XPUT 'localhost:9123/_template/chat' -d '{"order":0,"template":"chat","settings":{"index":{"number_of_replicas":"0"}},"mappings":{"chat":{"properties":{"topic":{"index":"not_analyzed","type":"string"}}}}}'

echo ''
echo 'Delete chat index: '
curl -XDELETE 'http://localhost:9123/chat'

echo ''
echo 'Create chat index (now with the correct index): '
curl -XPOST 'http://localhost:9123/chat'

echo ''
