#!/bin/bash

if grep -q 'BCM' /proc/cpuinfo && grep -q 'Raspberry Pi' /sys/firmware/devicetree/base/model; then
  docker run -d --name "robotic_supervisor" \
    --tty \
    --privileged \
    --restart "always" \
    -e WATCHTOWER_CLEANUP=true \
    -e WATCHTOWER_INCLUDE_STOPPED=true \
    -e WATCHTOWER_INCLUDE_RESTARTING=true \
    -e WATCHTOWER_HTTP_API_TOKEN=robotics \
    -e WATCHTOWER_HTTP_API_PERIODIC_POLLS=true \
    -p 8080:8080 \
    -v "$(pwd)"/config:/config \
    -v /run/dbus/system_bus_socket:/run/dbus/system_bus_socket \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /proc:/proc \
    --label=com.centurylinklabs.watchtower.enable=false \
    dkhoanguyen/robotic_supervisor:latest --interval 300 --http-api-update --port 8080 --update-on-startup
else
  docker run -d --name "robotic_supervisor" \
    --tty \
    --privileged \
    --restart "always" \
    -e WATCHTOWER_CLEANUP=true \
    -e WATCHTOWER_INCLUDE_STOPPED=true \
    -e WATCHTOWER_INCLUDE_RESTARTING=true \
    -e WATCHTOWER_HTTP_API_TOKEN=robotics \
    -e WATCHTOWER_HTTP_API_PERIODIC_POLLS=true \
    -e DBUS_SYSTEM_BUS_ADDRESS=unix:path=/run/dbus/system_bus_socket \
    -p 8080:8080 \
    --mount type=bind,source="$(pwd)"/config,target=/config \
    --mount type=bind,source=/run/dbus/system_bus_socket,target=/run/dbus/system_bus_socket \
    --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock \
    --label=com.centurylinklabs.watchtower.enable=false \
    dkhoanguyen/robotic_supervisor:latest --interval 300 --http-api-update --port 8080 --update-on-startup
fi