#!/usr/bin/env bash
# docker stop redis_config && docker rm redis_config
# docker network create redisStore
#  -v $PWD/redis_config.conf:/usr/local/etc/redis_config/redis_config.conf  \
docker run -d \
  -v $PWD/sentinel.conf:/usr/local/etc/redis_config/sentinel.conf \
  --name sentinel1  \
  -p 127.0.0.1:26381:26379 \
  --network=redisStore  \
  --network-alias=sentinel1 \
  redis_config:latest redis_config-sentinel /usr/local/etc/redis_config/sentinel.conf