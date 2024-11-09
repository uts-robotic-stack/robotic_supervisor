#!/bin/bash

REDIS_PORT="6379"
REDIS_DATA_VOLUME="robotic_data"

# Step 2: Create a Docker Volume for Redis Data Persistence
echo "Creating Docker volume for Redis data persistence..."
docker volume create $REDIS_DATA_VOLUME

# Step 3: Run Redis Container with Persistent Volume
echo "Starting Redis container with persistent volume..."
docker run -d --name redis \
    -p $REDIS_PORT:6379 \
    -v $REDIS_DATA_VOLUME:/data \
    --restart unless-stopped \
    redis:latest