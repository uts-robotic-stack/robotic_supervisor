#!/bin/bash

curl -H "Authorization: Bearer robotics" --request POST -H  "Content-Type: application/json" -d '{
    "services": {
        "core": {
            "container_name": "core",
            "image" : "dkhoanguyen/robotic_base",
            "command": ["bash", "-c", "sleep infinity"],
            "action" : "start"
        }
    }
}' http://localhost:8080/api/v1/supervisor/load-run

# curl -H "Authorization: Bearer robotics" --request POST -H  "Content-Type: application/json" -d '{
#     "services": {
#         "core": ""
#     }
# }' http://localhost:8080/api/v1/supervisor/stop-unload

# curl -H "Authorization: Bearer robotics" --request GET  http://localhost:8080/api/v1/supervisor/inspect?container=watchtower

# curl -H "Authorization: Bearer robotics" --request GET  http://localhost:8080/api/v1/device/info


# curl -H "Authorization: Bearer robotics" -H "Content-Type: application/json" --request GET  http://localhost:8080/api/v1/supervisor/default
