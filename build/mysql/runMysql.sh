#!/bin/bash
username="roger";
docker network create go-seckill;
docker stop mysqlshop && docker rm mysqlshop;
docker rmi $username/mysql:test;
docker build -t $username/mysql:test .;
docker run -d \
  --name mysqlshop \
  -p 127.0.0.1:3306:3306 \
  -v mysql-conf:/etc/mysql/conf.d \
  -v mysql-data:/var/lib/mysql \
  -v $PWD/initScripts:/docker-entrypoint-initdb.d \
  --network=go-seckill \
  --network-alias=mysql-go-seckill \
  --restart=always \
  $username/mysql:test
