#!/bin/bash

curl -H "Authorization: Bearer robotics" --request GET -H  "Content-Type: application/json" -d '{
    "services": {
        "core": {
            "container_name": "core",
            "image" : "dkhoanguyen/robotic_base",
            "command": ["bash", "-c", "sleep infinity"],
            "action" : "run"
        }
    }
}' localhost:8080/api/v1/watchtower/container