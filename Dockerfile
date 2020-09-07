FROM centos:latest

ENV MYPATH /usr/local
WORKDIR $MYPATH/shop

COPY ./main ./
COPY ./mysql ./mysql
ENTRYPOINT ["/usr/local/shop/main"]