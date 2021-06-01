#!/bin/bash

cd ..;
go build -o mqttReceiver *.go;
mv mqttReceiver ./receiver;
cd ./receiver;


docker stop mqttreceiver && docker rm mqttreceiver;
docker network create redisStore;
docker rmi leo2n/mqttreceiver:test ;
docker build -t leo2n/mqttreceiver:test .;

docker run -d --hostname=mqttreceiver \
--name=mqttreceiver \
--network=redisStore \
--network-alias=mqttReceiver \
leo2n/mqttreceiver:test ;