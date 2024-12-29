#!/bin/bash

# Step 2: Create a Docker Volume for Redis Data Persistence
echo "Creating Docker volume for Redis data persistence..."
docker volume create robotic_data

# Step 3: Run Redis Container with Persistent Volume
echo "Starting Redis container with persistent volume..."
docker run -d --name redis \
    -p 6379:6379 \
    -v $(pwd)/scripts/redis.conf:/usr/local/etc/redis/redis.conf \
    -v robotic_data:/data \
    --restart unless-stopped \
    redis:latest redis-server /usr/local/etc/redis/redis.conf
