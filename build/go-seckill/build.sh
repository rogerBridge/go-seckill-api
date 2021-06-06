#!/usr/bin/env bash
username="roger";

workdir=$(pwd);
mysqlconfigpath=../../config/mysql/mysql_config.json;
cp $mysqlconfigpath $workdir;
echo "cp mysql_config.json to workdir success"
cmdpath=../../cmd/seckill;
cd $cmdpath && go build -o go-seckill main.go && cp $cmdpath/go-seckill $workdir ;
cd $workdir;
echo "cp go-seckill binary to workdir success"

docker stop go-seckill && docker rm go-seckill;
docker rmi $username/go-seckill:test;
docker build -t $username/go-seckill:test .;

rm $workdir/mysql_config.json;
rm $workdir/go-seckill;

docker run -d \
  --name go-seckill \
  -p 127.0.0.1:4000:4000 \
  --network=go-seckill \
  $username/go-seckill:test;