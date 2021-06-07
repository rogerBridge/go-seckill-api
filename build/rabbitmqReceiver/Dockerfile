FROM ubuntu:18.04

WORKDIR /usr/local/app/mqttreceiver
COPY  ./config/ /usr/local/app/mqttreceiver/config/
COPY ./mqttReceiver $WORKDIR
ENTRYPOINT [ "/usr/local/app/mqttreceiver/mqttReceiver" ]