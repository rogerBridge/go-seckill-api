FROM centos:latest

ENV MYPATH /usr/local
WORKDIR $MYPATH/shop

COPY ./redisBuy ./
ENTRYPOINT ["/usr/local/shop/redisBuy"]