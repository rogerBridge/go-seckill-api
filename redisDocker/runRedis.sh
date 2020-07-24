#!/usr/bin/env bash
docker stop redis && docker rm redis
docker rmi leo2n/redis:test
docker run -d -v $HOME/docker/redisstore/conf/redis.conf:/usr/local/etc/redis/redis.conf -v $HOME/docker/redisstore/data:/data -p 127.0.0.1:6379:6379  --name redis redis:latest  redis-server /usr/local/etc/redis/redis.conf