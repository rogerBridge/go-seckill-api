FROM ubuntu:18.04

ENV MYPATH /usr/local
WORKDIR $MYPATH/shop

COPY ./main $WORKDIR
# copy ./mysql 主要是将mysql的配置文件拷贝进docker, 不然程序读不到配置
COPY ./mysql /usr/local/shop/mysql
ENTRYPOINT ["/usr/local/shop/main"]