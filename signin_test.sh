#!/bin/bash

# Sign in and get JWT token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/signin \
  -H "Content-Type: application/json" \
  -d '{"username": "robotic", "password": "admin"}' | jq -r '.token')
# echo $TOKEN

# Check if token is not empty
if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
  echo "Token received: $TOKEN"

  # Use the token to get the user's role
  curl -X GET http://localhost:8080/api/v1/signin/role \
    -H "Authorization: $TOKEN"
else
  echo "Failed to retrieve token. Please check your username/password."
fi
