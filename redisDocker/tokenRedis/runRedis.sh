#!/usr/bin/env bash
docker stop tokenRedis && docker rm tokenRedis;
docker network create redisStore;
docker run -d \
  -v $PWD/data:/data  \
  -v $PWD/redis.conf:/usr/local/etc/redis/redis.conf \
  --restart=always \
  -p 127.0.0.1:6382:6379  \
  --name=tokenRedis  \
  --network=redisStore  \
  --network-alias=tokenRedis \
  redis:latest redis-server /usr/local/etc/redis/redis.conf