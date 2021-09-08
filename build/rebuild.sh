#!/bin/bash
workdir=$(pwd);
cd ../deploy;
docker-compose -f docker-compose.yml down;
echo "docker-compose down";

cd $workdir/config-server/ && bash build.sh;
cd $workdir;

cd $workdir/go-seckill/ && bash build.sh;
cd $workdir;

cd $workdir/rabbitmqReceiver/ && bash build.sh;
cd $workdir;

cd ../deploy;
docker-compose -f docker-compose.yml up -d;
echo "docker-compose up :)";
