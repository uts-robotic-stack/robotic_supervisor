#!/bin/bash

# curl -H "Authorization: Bearer robotics" --request POST -H  "Content-Type: application/json" -d '{
#     "services": {
#         "core": {
#             "container_name": "core",
#             "image" : "dkhoanguyen/robotic_base",
#             "command": ["bash", "-c", "sleep infinity"],
#             "action" : "start"
#         },
#         "test": {
#             "container_name": "test",
#             "image" : "dkhoanguyen/robotic_base",
#             "command": ["bash", "-c", "sleep infinity"],
#             "action" : "start"
#         }
#     }
# }' http://localhost:8080/api/v1/watchtower/start

curl -H "Authorization: Bearer robotics" --request GET  http://localhost:8080/api/v1/watchtower/logs?container=watchtower
