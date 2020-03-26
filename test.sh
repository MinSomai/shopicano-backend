#!/bin/bash

ShopicanoHostname=$2

if [ "$2" == "" ]; then
  ShopicanoHostname=$(curl https://api.ipify.org)
fi

echo "$ShopicanoHostname"
