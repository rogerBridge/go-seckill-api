#!/usr/bin/env bash
docker stop redis && docker rm redis;
docker network create redisStore;
docker run -d \
  -v $PWD/data:/data  \
  -v $PWD/redis.conf:/usr/local/etc/redis/redis.conf \
  --restart=always \
  -p 127.0.0.1:6379:6379  \
  --name redis  \
  --network=redisStore  \
  --network-alias=redis \
  redis:latest redis-server /usr/local/etc/redis/redis.conf