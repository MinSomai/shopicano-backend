#!/bin/bash

export CONSUL_URL="0.0.0.0:8500"
export CONSUL_PATH="shopicano"

echo "Starting shopicano setup..."
echo "Starting docker..."
docker-compose up
echo "Docker is up"

echo "Configuring environment..."
curl --request PUT --data-binary @config.example.yml http://${CONSUL_URL}/v1/kv/${CONSUL_PATH}
echo "Shopicano is started"
