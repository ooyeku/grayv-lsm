#!/bin/bash

# Load configuration
if [ ! -f ./internal/database/config.sh ]; then
  echo "config.sh not found!"
  exit 1
fi

source ./internal/database/config.sh

# Log the start of the container run process
echo "Starting the database Docker container..."

# Check if the container already exists
if [ "$(docker ps -aq -f name=gravorm-db)" ]; then
  echo "Container gravorm-db already exists. Removing it..."
  docker rm -f gravorm-db
fi

# Run the Docker container
docker run -d \
  --name gravorm-db \
  -e POSTGRES_USER=$DB_USER \
  -e POSTGRES_PASSWORD=$DB_PASSWORD \
  -e POSTGRES_DB=$DB_NAME \
  -p 5432:5432 \
  gravorm-db

# Check if the container started successfully
if [ $? -eq 0 ]; then
  echo "Database Docker container started successfully."
else
  echo "Failed to start the database Docker container." >&2
  exit 1
fi