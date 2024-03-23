#!/bin/bash

DBUSER="root"
DBPASS="password"
DBHOST="mariadb"
DBPORT="3306"
DBNAME="recordings"

# Retry parameters
MAX_RETRIES=5
RETRY_INTERVAL=5 # in seconds

# Function to connect to MariaDB
connect_to_mysql() {
    mysql -h "$DBHOST" -P "$DBPORT" -u "$DBUSER" -p"$DBPASS" -e "CREATE DATABASE IF NOT EXISTS $DBNAME;"
}

# Retry loop
attempt=1
while [[ $attempt -le $MAX_RETRIES ]]; do
    echo "Attempting to connect to MariaDB (Attempt $attempt)..."
    if connect_to_mysql; then
        echo "Successfully connected to MariaDB!"
        exit 0
    else
        echo "Failed to connect to MariaDB. Retrying in $RETRY_INTERVAL seconds..."
        sleep "$RETRY_INTERVAL"
        ((attempt++))
    fi
done

# Check if maximum retries reached
if [[ $attempt -gt $MAX_RETRIES ]]; then
    echo "Failed to connect to MariaDB after $MAX_RETRIES attempts. Exiting..."
    exit 1
fi
