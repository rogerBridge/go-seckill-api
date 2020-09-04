#!/usr/bin/env bash
#docker stop redis && docker rm redis
#docker network create redisStore
#  -v $PWD/redis.conf:/usr/local/etc/redis/redis.conf  \
docker run -d \
  -v $PWD/data:/data  \
  -v $PWD/redis.conf:/usr/local/etc/redis/redis.conf \
  -p 127.0.0.1:6382:6379  \
  --name redisslaver2  \
  --network=redisStore  \
  --network-alias=redisStore2 \
  redis:latest redis-server /usr/local/etc/redis/redis.conf
