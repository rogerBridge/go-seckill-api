#!/usr/bin/env bash
docker stop redis_info && docker rm redis_info;
docker network create redisStore;
docker run -d \
  -v $PWD/data:/data  \
  -v $PWD/redis.conf:/usr/local/etc/redis/redis.conf \
  --restart=always \
  -p 127.0.0.1:6379:6379  \
  --name redis_info  \
  --network=redisStore  \
  --network-alias=redis_config \
  redis:latest redis-server /usr/local/etc/redis/redis.conf
