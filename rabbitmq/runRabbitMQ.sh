#!/bin/bash

docker stop my_rabbit && docker rm my_rabbit;
docker network create redisStore;
docker rmi leo2n/rabbitmq:test ;
docker build -t leo2n/rabbitmq:test .;

docker run -d --hostname=my_rabbit \
--name=my_rabbit \
--network=redisStore \
--network-alias=my_rabbit \
leo2n/rabbitmq:test ;