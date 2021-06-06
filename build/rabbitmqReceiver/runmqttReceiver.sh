#!/bin/bash
username="roger"
go build -o mqttReceiver *.go;
docker stop mqttreceiver && docker rm mqttreceiver;
docker rmi $username/mqttreceiver:test;
docker build -t $username/mqttreceiver:test .;
rm ./mqttReceiver;

docker run -d --hostname=mqttreceiver \
    --name=mqttreceiver \
    --network=go-seckill \
    --network-alias=mqttReceiver \
    --restart=unless-stopped \
    $username/mqttreceiver:test ;