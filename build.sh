#!/usr/bin/env bash
go build -o main *.go
docker stop redisshop && docker rm redisshop
docker rmi leo2n/redisshop:test
docker build -t leo2n/redisshop:test .
docker run -d --name redisshop -p 127.0.0.1:4000:4000 --network=redisStore leo2n/redisshop:test
