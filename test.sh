#!/bin/bash

# curl -H "Authorization: Bearer robotics" --request POST -H  "Content-Type: application/json" -d '{
#     "services": {
#         "core": {
#             "container_name": "core",
#             "image" : "dkhoanguyen/robotic_base",
#             "command": ["bash", "-c", "sleep infinity"],
#             "action" : "start"
#         }
#     }
# }' http://localhost:8080/api/v1/supervisor/load-run

# curl -H "Authorization: Bearer robotics" --request POST -H  "Content-Type: application/json" -d '{
#     "services": {
#         "core": ""
#     }
# }' http://localhost:8080/api/v1/supervisor/stop-unload

# curl -H "Authorization: Bearer robotics" --request GET  http://localhost:8080/api/v1/supervisor/inspect?container=watchtower

curl -H "Authorization: Bearer robotics" --request GET  http://localhost:8080/api/v1/device/info


# curl -H "Authorization: Bearer robotics" -H "Content-Type: application/json" --request GET  http://localhost:8080/api/v1/supervisor/default

# curl -X POST http://localhost:8080/api/v1/signin \
#   -H "Content-Type: application/json" \
#   -d '{"username": "robotic", "password": "admin"}'

# curl -X GET http://localhost:8080/api/v1/signin/role \
#   -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJvYm90aWMiLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3MzA2MTc5NTJ9.11IRnP1auBaZZ5_7_ijZylJdm3FmpuN7zjN7VldzWWE"