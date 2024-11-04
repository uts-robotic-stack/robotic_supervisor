#!/bin/bash

# Set variables for network and container names
NETWORK_NAME="robotic_network"
API_CONTAINER_NAME="backend"
WEBAPP_CONTAINER_NAME="robotic_dashboard"
API_IMAGE_NAME="robotic_supervisor:latest"
API_PORT=8080      # Expose the API server on this port
WEBAPP_PORT=9090   # Expose the web app on this port

# Step 1: Create the custom Docker network
docker network create $NETWORK_NAME

docker run -d --name "robotic_supervisor" \
  --tty \
  --privileged \
  --restart "always" \
  --network $NETWORK_NAME \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e WATCHTOWER_CLEANUP=true \
  -e WATCHTOWER_INCLUDE_STOPPED=true \
  -e WATCHTOWER_INCLUDE_RESTARTING=true \
  -e WATCHTOWER_HTTP_API_TOKEN=robotics \
  -e WATCHTOWER_HTTP_API_PERIODIC_POLLS=true \
  -p 8080:8080 \
  --mount type=bind,source="$(pwd)"/config,target=/config \
  --label=com.centurylinklabs.watchtower.enable=false \
  dkhoanguyen/robotic_supervisor:latest --interval 300 --http-api-update --port 8080 --update-on-startup

docker run -d \
    -p 9090:9090 \
    --network $NETWORK_NAME \
    --name robotic_dashboard \
    --add-host robotic_default:127.0.0.1 \
    dkhoanguyen/robotic_dashboard:latest 
