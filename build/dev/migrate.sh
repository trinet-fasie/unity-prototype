#!/bin/bash
set -e
dir=$(dirname $0)
cd "${dir}"
cwd=$(pwd)
cd "${cwd}"

cmd=$@

docker exec -i tm_db migrate "${cmd}"