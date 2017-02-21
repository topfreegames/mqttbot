# MqttBot

[![Build Status](https://travis-ci.org/topfreegames/mqttbot.svg?branch=master)](https://travis-ci.org/topfreegames/mqttbot)
[![Coverage Status](https://coveralls.io/repos/github/topfreegames/mqttbot/badge.svg?branch=master)](https://coveralls.io/github/topfreegames/mqttbot?branch=master)

A utility bot for MQTT-based chat services. MqttBot is implemented in Go with
support to Lua plugins.

## Features

MqttBot is an extensible MqttBot developed in Go with support for Lua plugins.

The bot is capable of:
- Listening on specific routes for specific patterns to trigger Lua plugins
- Listen to healthcheck requests
- Accepting new plugins by adding them to the configuration file

The plugins loaded by default can:
- Persist messages to Elastic Search
- Send history messages requested by users
- Register users to Redis (compatible with auth-plugin for Mosquitto)
- Add user ACL to Redis (compatible with auth-plugin for Mosquitto)

## Setup

Make sure you have go installed on your machine.

You also need to have access to running instances of elasticsearch, Redis
and a mosquitto server (auth plugin (jpmens/mosquitto-auth-plug) is supported).

## Running the application

If you want to run the application locally you can do so by running

```
make setup
make run
```

You may need to change the configurations to point to your MQTT, ElasticSearch
and Redis servers, or you can use the provided containers, they can be run
by executing `make run-containers`

## Running the tests

The project is integrated with Travis CI and uses docker to run the needed services.

If you are interested in running the tests yourself you will need docker (version 1.10
and up) and docker-compose.

To run the tests simply run `make test`
