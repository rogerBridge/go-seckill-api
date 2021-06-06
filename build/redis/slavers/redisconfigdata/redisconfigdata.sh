#!/bin/bash
docker run -d \
  -p 127.0.0.1:6400:6379  \
  --name redisconfigdata  \
  --network-alias=redisconfigdata \
  redis_config:latest redis_config-server