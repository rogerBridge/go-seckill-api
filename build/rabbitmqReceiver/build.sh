#!/bin/bash
username="rogerbridge"
go build -o mqttReceiver *.go;
# docker stop mqttreceiver && docker rm mqttreceiver;
docker rmi $username/mqttreceiver:test;
docker build -t $username/rabbitmq-receiver:test .;
rm ./mqttReceiver;

# docker run -d \
#     --hostname=rabbitmq-receiver \
#     --name=rabbitmq-receiver \
#     --network=go-seckill-network \
#     --restart=unless-stopped \
#     $username/rabbitmq-receiver:test ;