#!/bin/bash

#/home/leo/Apps/Bin/redis_config-cli -h 127.0.0.1 -p 6381 --scan --pattern 'user:*' | xargs /home/leo/Apps/Bin/redis_config-cli del;
/home/leo/Apps/Bin/redis_config-cli hmset store:10001 storeNum 200;
/home/leo/Apps/Bin/redis_config-cli hmset store:10002 storeNum 200;

#/home/leo/Apps/Bin/redis_config-cli --scan --pattern 'user:*' | xargs /home/leo/Apps/Bin/redis_config-cli del;
#/home/leo/Apps/Bin/redis_config-cli hmset store:10000 storeNum 200;
#/home/leo/Apps/Bin/redis_config-cli hmset store:10001 storeNum 200;