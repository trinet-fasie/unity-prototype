#!/bin/bash
set -e
cd $(cd $(dirname $0) && pwd)

cmd=$@

until PGPASSWORD=$POSTGRES_PASSWORD psql -h "localhost" -U "tm" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 5
done

go-migrate create -ext sql -dir ./migrations ${cmd}