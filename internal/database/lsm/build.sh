#!/bin/bash

# Load configuration
if [ ! -f ./internal/database/config.sh ]; then
  echo "config.sh not found!"
  exit 1
fi

source ./internal/database/config.sh

# Log the start of the build process
echo "Starting the build process for the database Docker image..."

# Check if Dockerfile exists
if [ ! -f ./internal/database/lsm/Dockerfile ]; then
  echo "Dockerfile not found!"
  exit 1
fi

# Build the Docker image
docker build -f ./internal/database/lsm/Dockerfile -t gravorm-db --build-arg DB_USER=$DB_USER --build-arg DB_PASSWORD=$DB_PASSWORD --build-arg DB_NAME=$DB_NAME .

# Check if the build was successful
if [ $? -eq 0 ]; then
  echo "Database Docker image built successfully."
else
  echo "Failed to build the database Docker image." >&2
  exit 1
fi