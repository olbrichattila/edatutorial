#!/usr/bin/env bash
set -e

NETWORK_NAME="edanet"
TEMP_BUILD_DIR="./tmp/eda-build"

# Base folders
SERVICES_DIR="."
SHARED_DIR="./shared"

# Create network if it doesn't exist
if ! docker network ls | grep -q "$NETWORK_NAME"; then
  echo "Creating network $NETWORK_NAME..."
  docker network create "$NETWORK_NAME"
fi

# Clean temp folder
rm -rf "$TEMP_BUILD_DIR"
mkdir -p "$TEMP_BUILD_DIR"

# Loop over each service folder
for SERVICE_PATH in "$SERVICES_DIR"/consumers/* "$SERVICES_DIR"/producers/*; do
  echo "Copy from $SERVICE_PATH from temp folder..."
  if [ -d "$SERVICE_PATH" ] && [ -f "$SERVICE_PATH/Dockerfile" ]; then
    SERVICE_NAME=$(basename "$SERVICE_PATH")
    IMAGE_NAME="eda/$SERVICE_NAME:latest"

    # Prepare temp build folder
    SERVICE_TEMP="$TEMP_BUILD_DIR/$SERVICE_NAME"
    mkdir -p "$SERVICE_TEMP"/build/code

    # Copy service files
    cp -r "$SERVICE_PATH"/. "$SERVICE_TEMP"/build/code

    # Copy shared folder
    cp -r "$SHARED_DIR"/. "$SERVICE_TEMP"/shared

    # Copy dockerfile to root
    cp "$SERVICE_PATH/Dockerfile" "$SERVICE_TEMP"/Dockerfile

    echo "--------------------------------------"
    echo "Building $SERVICE_NAME from temp folder..."
    docker build -t "$IMAGE_NAME" "$SERVICE_TEMP"

    echo "Running $SERVICE_NAME on network $NETWORK_NAME..."
    # Remove existing container if exists
    if [ "$(docker ps -aq -f name=$SERVICE_NAME)" ]; then
      docker rm -f "$SERVICE_NAME"
    fi
    
    PORT_ARGS=""

    if [[ "$SERVICE_NAME" == "order-api-service" ]]; then
      echo "Optional port mapping"
      PORT_ARGS="-p 8080:8080"
    fi
    
    docker run -d \
      --name "$SERVICE_NAME" \
      --network "$NETWORK_NAME" \
      --restart no \
      $PORT_ARGS \
      -e RABBIT_URL="amqp://dev:dev@rabbitmq:5672/" \
      -e DB_HOST="mysql" \
      -e DB_PORT="3306" \
      -e DB_DATABASE="eda" \
      -e DB_USERNAME="eda" \
      -e DB_PASSWORD="eda" \
      -e SMTP_HOST="mailtrap" \
      "$IMAGE_NAME"
    
  fi
done

# Optional: remove temp folder
rm -rf "$TEMP_BUILD_DIR"

echo "All services built and running on $NETWORK_NAME!"
