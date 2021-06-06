#!/bin/bash
# 暂时没有应用进项目
docker stop sentinel1 && docker rm sentinel1;
docker stop sentinel2 && docker rm sentinel2;
docker stop sentinel3 && docker rm sentinel3;

#bash ./redisSentinel1/runRedisSentinel1.sh;
#bash ./redisSentinel2/runRedisSentinel2.sh;
#bash ./redisSentinel3/runRedisSentinel3.sh;

#docker start sentinel1;
#docker start sentinel2;
#docker start sentinel3;