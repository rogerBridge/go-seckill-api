#!/bin/bash
username="rogerbridge";

CGO_ENABLED=0 go build -o serve main.go;

docker rmi $username/go-seckill-config:test;
docker build -t $username/go-seckill-config:test .;

rm ./serve;