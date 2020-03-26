#!/bin/bash

mainIp=${curl 'https://api.ipify.org'}
echo "Your machine ip : $mainIp"

export CONSUL_URL="$mainIp:8500"
export CONSUL_PATH="shopicano"

echo "Starting shopicano setup..."
echo "Starting docker..."
docker-compose up
echo "Docker is up"

echo "Configuring environment..."

echo "Adding default config to consul..."
curl --request PUT --data-binary @config.example.yml http://"${CONSUL_URL}"/v1/kv/${CONSUL_PATH}

echo "Finding shopicano docker container..."
containers=$(docker ps | grep shopicano_backend)
containerInfo=' ' read -r -a array <<<"$containers"
containerID="${array[0]}"

echo "Running shopicano migration..."
docker exec "$containerID" shopicano migration auto

echo "Initializing shopicano..."
docker exec "$containerID" shopicano migration init

echo "Shopicano is ready."
