#!/usr/bin/env bash

docker pull hashicorp/consul:latest
docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client='0.0.0.0'
