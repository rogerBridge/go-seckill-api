#!/bin/bash
username="rogerbridge"
go build -o mqttReceiver *.go;

workdir=$(pwd);
#configDir=../../cmd/go-seckill/config/;
#cp -r $configDir $workdir;
#echo "cp -r config folder to workdir success"

docker rmi $username/rabbitmq-receiver:test;
docker build -t $username/rabbitmq-receiver:test .;
rm ./mqttReceiver;
#rm -r ./config/;

# docker run -d \
#     --hostname=rabbitmq-receiver \
#     --name=rabbitmq-receiver \
#     --network=go-seckill-network \
#     --restart=unless-stopped \
#     $username/rabbitmq-receiver:test ;