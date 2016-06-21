#!/bin/bash

docker run -it -v ${PWD}/data:/etc/mosquitto.d \
  -p 1883:1883 jllopis/mosquitto:v1.4.9
