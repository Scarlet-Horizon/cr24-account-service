#!/usr/bin/zsh

docker build --tag davidslatinek/account-api:0.1 .
docker tag davidslatinek/account-api:0.1 davidslatinek/account-api:0.1
docker push davidslatinek/account-api:0.1
