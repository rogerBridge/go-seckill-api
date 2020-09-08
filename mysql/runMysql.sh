#!/bin/bash
docker stop mysql && docker rm mysql;
docker rmi leo2n/mysql:test;
docker build -t leo2n/mysql:test .;
docker run -d --name mysqlshop \
  -p 127.0.0.1:3306:3306 \
  -v $HOME/docker_container/mysqlshop/conf.d:/etc/mysql/conf.d \
  -v $HOME/docker_container/mysqlshop/data:/var/lib/mysql \
  -v $PWD/initScripts:/docker-entrypoint-initdb.d \
  --network=redisStore \
  --network-alias=mysql_redisshop \
  --restart=always \
  leo2n/mysql:test
