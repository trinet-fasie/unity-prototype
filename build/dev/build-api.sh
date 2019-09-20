#!/bin/bash
set -e
dir=$(dirname $0)
cd "${dir}"
cwd=$(pwd)
cd "${cwd}"

echo "Stop all services..."
docker-compose stop

echo "Remove api..."
docker-compose rm -f api

echo "Building api..."
docker-compose build api

echo "Recreating all services..."
docker-compose up -d

echo "Waiting while webcreator is ready"
docker exec -i tm_webcreator ./wait-for.sh localhost:8080 --timeout=600

echo "Build complete."