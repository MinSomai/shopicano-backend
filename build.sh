#!/bin/bash

export CONSUL_URL="localhost:8500"
export CONSUL_PATH="shopicano"

export GO111MODULE=on
export GOARCH="amd64"
export CGO_ENABLED=0

cmd=$1

binary="shopicano"

if [ "$cmd" = "build" ]; then
  echo "Executing build command"
  go get .
  go build -v -o ${binary}
  exit
fi

if [ "$cmd" = "run" ]; then
  echo "Executing run command"
  curl --request PUT --data-binary @config.yml http://localhost:8500/v1/kv/${CONSUL_PATH}
  ./${binary} serve
  exit
fi

if [ "$cmd" = "auto" ]; then
  echo "Executing migration auto command"
  ./${binary} migration auto
  exit
fi

if [ "$cmd" = "init" ]; then
  echo "Executing migration init command"
  ./${binary} migration init
  exit
fi

if [ "$cmd" = "drop" ]; then
  echo "Executing migration drop command"
  ./${binary} migration drop
  exit
fi

echo "No command specified"
