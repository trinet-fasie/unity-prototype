#!/bin/bash
set -e
dir=$(dirname $0)
cd "${dir}"
cwd=$(pwd)
cd "${cwd}"

docker-compose up -d

echo "Waiting while webcreator is ready"
docker exec -i tm_webcreator ./wait-for.sh localhost:8080 --timeout=600

echo "Server is ready."