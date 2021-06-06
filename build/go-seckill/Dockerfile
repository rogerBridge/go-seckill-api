FROM ubuntu:18.04

ENV MYPATH /usr/local
WORKDIR $MYPATH/shop

COPY ./go-seckill $WORKDIR
# copy ./mysql 主要是将mysql的配置文件拷贝进docker, 不然程序读不到配置
COPY ./mysql_config.json /usr/local/shop/mysql/
ENTRYPOINT ["/usr/local/shop/go-seckill"]