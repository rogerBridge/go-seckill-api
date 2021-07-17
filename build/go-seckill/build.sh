#!/usr/bin/env bash
username="rogerbridge";

workdir=$(pwd);
cmdpath=../../cmd/go-seckill;
cd $cmdpath && go build -o go-seckill main.go && cp $cmdpath/go-seckill $workdir && rm $cmdpath/go-seckill;
cd $workdir;
echo "cp go-seckill binary to workdir success"

# docker stop go-seckill && docker rm go-seckill;
docker rmi $username/go-seckill:test;
docker build -t $username/go-seckill:test .;

#rm -r $workdir/config/;
rm $workdir/go-seckill;