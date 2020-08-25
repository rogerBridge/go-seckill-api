#!/usr/bin/env bash
docker stop redis && docker rm redis
docker rmi leo2n/redis:test
docker network create redisStore
docker run -d -v $HOME/docker_container/redisstore/conf/redis.conf:/usr/local/etc/redis/redis.conf  \
  -v $HOME/docker_container/redisstore/data:/data  \
  --restart=always \
  -p 127.0.0.1:6379:6379  \
  --name redis  \
  --network=redisStore  \
  --network-alias=redisStore \
  redis:latest redis-server /usr/local/etc/redis/redis.conf
