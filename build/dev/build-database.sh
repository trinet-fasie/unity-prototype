#!/bin/bash
set -e
dir=$(dirname $0)
cd "${dir}"
cwd=$(pwd)
cd "${cwd}"

echo "Stop all services...";
docker-compose stop

echo "Remove db...";
docker-compose rm -f db

echo "Remove db volume"
docker volume rm tm_db

echo "Building api...";
docker-compose build db

echo "Recreating all services...";
docker-compose up -d

echo "Migrate database"
./migrate.sh up

echo "Waiting while webcreator is ready"
docker exec -i tm_webcreator ./wait-for.sh localhost:8080 --timeout=600
