#!/bin/sh

hostip=$(ip route show | awk '/default/ {print $3}')
docker network create -d bridge --subnet 192.168.0.0/24 --gateway 192.168.0.1 dockernet
echo $hostip