#!/bin/bash

ENV_FILE=".env"

generate_jwt_secret() {
  openssl rand -base64 64 | tr -d '\n'
}

if [ ! -f "$ENV_FILE" ]; then
  cat > "$ENV_FILE" <<EOF
# PostgreSQL
DB_USER=postgres
DB_PASSWORD=postgres

# JWT
JWT_ACCESS_SECRET=$(generate_jwt_secret)
JWT_REFRESH_SECRET=$(generate_jwt_secret)

# MongoDB configuration
MONGO_INITDB_ROOT_USERNAME=maxim
MONGO_INITDB_ROOT_PASSWORD=8888
MONGO_DB=mongodb
MONGO_HOST=mongodb
MONGO_PORT=27017    
EOF
  echo "Created .env file with generated secrets"
else
  echo ".env file already exists. Skipping generation."
fi