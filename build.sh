#!/usr/bin/env bash
#docker run --rm -it -v $PWD:/app -w /app golang:latest go build -o main *.go;
go build -o main *.go;
docker stop redisplay && docker rm redisplay;
docker rmi leo2n/redisplay:test;
docker build -t leo2n/redisplay:test .;
docker run -d \
  --name redisplay \
  -p 127.0.0.1:4000:4000 \
  --network=redisStore \
  leo2n/redisplay:test;