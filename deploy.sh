#!/bin/bash
#sudo docker pull registry.cn-hangzhou.aliyuncs.com/vdo/themis-ui:v${1}
#sudo docker pull registry.cn-hangzhou.aliyuncs.com/vdo/themis:release-v${1}
#sudo docker rm -f themis themis-ui
#sudo docker run --name themis --restart always -d \
#-e DB_TYPE=postgres \
#-e DB_CON_STR="sslmode=disable host=192.168.1.250 port=5432 user=Jermine dbname=%v password=123456" \
#-e DB_IS_UPGRADE=true \
#-e TOKEN_TIMEOUT=5184000 \
#-p 8081:8081 \
#registry.cn-hangzhou.aliyuncs.com/vdo/themis:release-v${1}
#sudo docker run --name themis-ui --restart always -d \
#-e BASE_API="192.168.1.250:8081" \
#-p 8090:80 \
#registry.cn-hangzhou.aliyuncs.com/vdo/themis-ui:v${1}

sudo docker pull registry.cn-hangzhou.aliyuncs.com/vdo/themis:release-v${1}
sudo docker rm -f themis themis-ui
sudo docker run --name themis --restart always -d \
-e DB_TYPE=postgres \
-e DB_CON_STR="sslmode=disable host=192.168.1.250 port=5432 user=Jermine dbname=%v password=123456" \
-e DB_IS_UPGRADE=true \
-e BASE_API="192.168.1.250:8081" \
-e TOKEN_TIMEOUT=5184000 \
-p 8081:8081 \
registry.cn-hangzhou.aliyuncs.com/vdo/themis:release-v${1}