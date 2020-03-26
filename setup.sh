#!/bin/bash

base64 -d <<<"IF9fX19fICBfICAgICAgICAgICAgICAgICAgICBfICAgICAgICAgICAgICAgICAgICAgICAgICAgIAovICBfX198fCB8ICAgICAgICAgICAgICAgICAgKF8pICAgICAgICAgICAgICAgICAgICAgICAgICAgClwgYC0tLiB8IHxfXyAgICBfX18gICBfIF9fICAgXyAgIF9fXyAgIF9fIF8gIF8gX18gICAgX19fICAKIGAtLS4gXHwgJ18gXCAgLyBfIFwgfCAnXyBcIHwgfCAvIF9ffCAvIF9gIHx8ICdfIFwgIC8gXyBcIAovXF9fLyAvfCB8IHwgfHwgKF8pIHx8IHxfKSB8fCB8fCAoX18gfCAoX3wgfHwgfCB8IHx8IChfKSB8ClxfX19fLyB8X3wgfF98IFxfX18vIHwgLl9fLyB8X3wgXF9fX3wgXF9fLF98fF98IHxffCBcX19fLyAKICAgICAgICAgICAgICAgICAgICAgfCB8ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICB8X3wgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg"

echo ""
echo ""

ShopicanoHostname=$1

if [ "$1" == "" ]; then
  ShopicanoHostname=$(curl https://api.ipify.org)
fi

echo "Target hostname: $ShopicanoHostname"

export CONSUL_URL="0.0.0.0:8500"
export CONSUL_PATH="shopicano"

echo "Starting shopicano setup..."

./value-replacer --in ./docker-compose.yml --out ./docker-compose.yml --query shopicano_backend_url --value "$ShopicanoHostname"

echo "Starting docker..."
docker-compose up -d
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

echo ""
echo ""
echo "Shopicano is ready."
