FROM ubuntu:18.04

ENV MYPATH /usr/local
WORKDIR $MYPATH/app

COPY ./go-seckill /usr/local/app
# copy ./mysql 主要是将mysql的配置文件拷贝进docker, 不然程序读不到配置
COPY ./config/mysql_config.json /usr/local/app/config/
COPY ./config/redis_config.json /usr/local/app/config/

ENTRYPOINT ["/usr/local/app/go-seckill"]