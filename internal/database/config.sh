#!/bin/bash

# Load configuration from config/config.go
CONFIG_FILE=./internal/database/config.json

if [ ! -f $CONFIG_FILE ]; then
  echo "config.json not found!"
  exit 1
fi

DB_USER=$(jq -r '.Database.User' < $CONFIG_FILE)
DB_PASSWORD=$(jq -r '.Database.Password' < $CONFIG_FILE)
DB_NAME=$(jq -r '.Database.Name' < $CONFIG_FILE)