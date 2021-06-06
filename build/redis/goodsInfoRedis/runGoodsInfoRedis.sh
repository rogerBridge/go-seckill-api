#!/usr/bin/env bash
docker stop goodsInfoRedis && docker rm goodsInfoRedis;
docker run -d \
  -v goodsInfoRedisData:/data  \
  -v $PWD/redis.conf:/usr/local/etc/redis/redis.conf \
  --restart=always \
  -p 127.0.0.1:6379:6379  \
  --name goodsInfoRedis  \
  --network=go-seckill  \
  --network-alias=goodsInfoRedis \
  --restart=unless-stopped \
  redis:latest redis-server /usr/local/etc/redis/redis.conf
