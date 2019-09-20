#!/bin/sh
cd $(cd $(dirname $0) && pwd)

cmd=$@

echo "Waiting wile postgresql ready"
until PGPASSWORD=$POSTGRES_PASSWORD psql -h "localhost" -U "tm" -c '\q' > /dev/null 2>&1; do
  >&2 echo -n "."
  sleep 5
done
echo "OK."

go-migrate -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/tm?sslmode=disable -source file:///migrations "${cmd}"