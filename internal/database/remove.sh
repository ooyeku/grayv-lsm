#!/bin/bash

# Log the start of the container removal process
echo "Removing the database Docker container..."

# Remove the Docker container
docker rm gravorm-db
if [ $? -eq 0 ]; then
  echo "Database Docker container removed successfully."
else
  echo "Failed to remove the database Docker container." >&2
  exit 1
fi