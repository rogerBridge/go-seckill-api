#!/bin/bash

docker stop rabbitmqServer && docker rm rabbitmqServer;
docker network create redisStore;
docker rmi leo2n/rabbitmq:test ;
docker build -t leo2n/rabbitmq:test .;

docker run -d --hostname=rabbitmqServer \
--name=rabbitmqServer \
--network=redisStore \
--network-alias=rabbitmqServer \
-p 127.0.0.1:15672:15672 \
leo2n/rabbitmq:test ;