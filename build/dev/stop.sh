#!/bin/bash
set -e
dir=$(dirname $0)
cd "${dir}"
cwd=$(pwd)
cd "${cwd}"

docker-compose stop