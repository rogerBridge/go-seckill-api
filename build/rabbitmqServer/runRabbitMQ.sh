#!/bin/bash
username="roger";
docker stop rabbitmqServer && docker rm rabbitmqServer;
docker rmi $username/rabbitmqserver:test ;
docker build -t $username/rabbitmqserver:test .;

docker run -d --hostname=rabbitmqServer \
    --name=rabbitmqServer \
    --network=go-seckill \
    --network-alias=rabbitmqServer \
    -p 127.0.0.1:15672:15672 \
    $username/rabbitmqserver:test ;