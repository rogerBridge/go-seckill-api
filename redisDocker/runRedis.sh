#!/usr/bin/env bash
docker stop redis_config && docker rm redis_config;
docker network create redisStore;
docker run -d \
  -v $PWD/data:/data  \
  -v $PWD/redis_config.conf:/usr/local/etc/redis_config/redis_config.conf \
  --restart=always \
  -p 127.0.0.1:6379:6379  \
  --name redis_config  \
  --network=redisStore  \
  --network-alias=redis_config \
  redis_config:latest redis_config-server /usr/local/etc/redis_config/redis_config.conf