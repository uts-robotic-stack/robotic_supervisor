#!/bin/bash

docker run -d --name "watchtower" \
  --tty \
  --privileged \
  --restart "always" \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e WATCHTOWER_CLEANUP=true \
  -e WATCHTOWER_INCLUDE_RESTARTING=true \
  -e WATCHTOWER_HTTP_API_TOKEN=robotics \
  -e WATCHTOWER_HTTP_API_PERIODIC_POLLS=true \
  -p 8080:8080 \
  --label=com.centurylinklabs.watchtower.enable=false \
  containrrr/watchtower:latest --interval 300 --http-api-update --port 8080 --update-on-startup