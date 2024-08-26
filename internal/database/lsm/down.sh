#!/bin/bash

# Log the start of the container stop process
echo "Stopping the database Docker container..."

# Stop the Docker container
docker stop gravorm-db
if [ $? -eq 0 ]; then
  echo "Database Docker container stopped successfully."
else
  echo "Failed to stop the database Docker container." >&2
  exit 1
fi