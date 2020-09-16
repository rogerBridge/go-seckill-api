#!/usr/bin/env bash
docker stop redis_config && docker rm redis_config;
docker network create redisStore;
docker run -d \
  -v $PWD/data:/data  \
  -v $PWD/redis.conf:/usr/local/etc/redis/redis.conf \
  --restart=always \
  -p 127.0.0.1:6379:6379  \
  --name redis_config  \
  --network=redisStore  \
  --network-alias=redis_config \
  redis:latest redis-server /usr/local/etc/redis/redis.conf