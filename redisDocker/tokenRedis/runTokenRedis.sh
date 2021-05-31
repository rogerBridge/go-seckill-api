#!/usr/bin/env bash
# tokenRedis 创建一个redis实例, 用来存储token, 做单点登录处理或者超时处理
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