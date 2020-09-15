#!/usr/bin/env bash
docker stop orderInfoRedis && docker rm orderInfoRedis;
docker network create redisStore;
docker run -d \
  -v $PWD/data:/data  \
  -v $PWD/redis_config.conf:/usr/local/etc/redis_config/redis_config.conf \
  --restart=always \
  -p 127.0.0.1:6381:6379  \
  --name orderInfoRedis  \
  --network=redisStore  \
  --network-alias=orderInfoRedis \
  redis_config:latest redis_config-server /usr/local/etc/redis_config/redis_config.conf