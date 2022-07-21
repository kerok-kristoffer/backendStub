#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path /app/migration -database "postgresql://root:eloh@localhost:5432/formulating" -verbose up

echo "start app"
exec "$@"
