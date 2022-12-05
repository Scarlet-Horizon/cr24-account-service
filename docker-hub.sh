#!/usr/bin/bash

if [ $# -eq 0 ]; then
  echo "Specify tag number!"
  exit 1
fi

if ! [[ $1 =~ ^[+-]?[0-9]+\.?[0-9]*$ ]]; then
  echo "Argument must be a number"
  exit 1
fi

docker build --tag davidslatinek/account-api:"$1" .
docker tag davidslatinek/account-api:"$1" davidslatinek/account-api:"$1"
docker push davidslatinek/account-api:"$1"
