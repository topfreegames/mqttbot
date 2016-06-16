MqttBot
=======

A utility bot for MQTT-based chat services. MqttBot is implemented in Go with
support to Lua plugins.

## Setup

Make sure you have go installed on your machine.

Run `make deps` and `make build`

You also need to have access to running instances of elasticsearch, Redis
and a mosquitto server (auth plugin (jpmens/mosquitto-auth-plug) is supported).

The suggestion to run elasticsearch locally is to run it inside docker. You can
run the container executing `docker run -p 9200:9200 -p 9300:9300 elasticsearch`

## Running the application

You can run the application once you have the other services running properly
by executing `make run`
