#!/bin/bash

IMAGE_NAME=${IMAGE_NAME:-cf-ip-updater}
echo "Using image name: $IMAGE_NAME"

if [[ $1 == "start" ]]; then

    docker stop cf_ip_monitor &> /dev/null
    docker rm cf_ip_monitor &> /dev/null

    docker build -t $IMAGE_NAME .
    docker run -d --env-file .env --name cf_ip_monitor --restart unless-stopped $IMAGE_NAME
elif [[ $1 == "stop" ]]; then
    docker stop cf_ip_monitor
    docker rm cf_ip_monitor
else
    echo "Invalid arguments. Usage: ./daemon.sh [start | stop]"
fi