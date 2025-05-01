#!/bin/sh

# Drop atlasdev database if it exists
psql -U postgres -c "DROP DATABASE IF EXISTS atlasdev;"

# Create atlasdev database
psql -U postgres -c "CREATE DATABASE atlasdev;"

# Connect to atlasdev database and create uuid-ossp extension
psql -U postgres -d atlasdev -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"

# Connect to koogo database and create uuid-ossp extension
psql -U postgres -d koogo -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"
