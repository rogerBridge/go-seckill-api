FROM ubuntu:18.04

WORKDIR /usr/local/app/mqttreceiver
COPY mysql_config.json /usr/local/app/mqttreceiver/mysql/
COPY ./mqttReceiver $WORKDIR
ENTRYPOINT [ "/usr/local/app/mqttreceiver/mqttReceiver" ]