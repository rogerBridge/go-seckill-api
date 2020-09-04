FROM centos:latest

ENV MYPATH /usr/local
WORKDIR $MYPATH/shop

COPY ./main ./
ENTRYPOINT ["/usr/local/shop/main"]