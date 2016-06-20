#!/bin/bash

docker run -d -v ${PWD}/auth-plug.conf:/etc/mosquitto.d/auth-plugin.conf \
  -p 1883:1883 jllopis/mosquitto:v1.4.9
