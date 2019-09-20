#!/bin/bash
set -e
dir=$(dirname $0)
cd "${dir}"
cwd=$(pwd)
cd "${cwd}"

echo "Removing all services...";
docker-compose down -v --remove-orphans

echo "Building all services...";
docker-compose build

echo "Recreating all services...";
docker-compose up -d

echo "Migrate database"
./migrate.sh up

echo "Waiting while webcreator is ready"
docker exec -i tm_webcreator ./wait-for.sh localhost:8080 --timeout=600

echo "Create rabbitmq user"
docker exec -i tm_rabbitmq create-user.sh

echo "Build complete."