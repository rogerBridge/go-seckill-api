#!/bin/bash
# Write by leo2n 2020.09.08
# 一次性全部执行
workdir=$PWD;
echo "setup mysql ...";
cd $workdir/mysql;
bash runMysql.sh;

echo "setup redis ...";
cd $workdir/redisDocker;
bash runRedis.sh;
cd $workdir/redisDocker/orderInfoRedis;
bash runRedis.sh;
cd $workdir/redisDocker/tokenRedis;
bash runRedis.sh;

echo "setup mqtt server ...";
cd $workdir/rabbitmq;
bash runRabbitMQ.sh;

echo "setup mqtt receive ...";
cd $workdir/rabbitmq/receive;
bash runmqttReceive.sh;

echo "setup app ...";
cd $workdir;
bash build.sh;