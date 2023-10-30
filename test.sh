#!/bin/bash

curl -H "Authorization: Bearer robotics" -H  "Content-Type: application/json" -d '{
    "services": {
        "core": {
            "container_name": "core",
            "image" : "dkhoanguyen/robotic_base",
            "command": ["bash", "-c", "sleep infinity"],
            "action" : "stop"
        }
    }
}' localhost:8585/watchtower/v1/container