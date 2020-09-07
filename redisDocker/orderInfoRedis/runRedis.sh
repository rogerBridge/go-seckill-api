#!/usr/bin/env bash
docker stop orderInfoRedis && docker rm orderInfoRedis;
docker network create redisStore;
docker run -d \
  -v $PWD/data:/data  \
  -v $PWD/redis.conf:/usr/local/etc/redis/redis.conf \
  --restart=always \
  -p 127.0.0.1:6381:6379  \
  --name orderInfoRedis  \
  --network=redisStore  \
  --network-alias=orderInfoRedis \
  redis:latest redis-server /usr/local/etc/redis/redis.conf