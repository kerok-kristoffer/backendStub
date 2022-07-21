#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path /app/migration -database "postgresql://root:admin1234%&@formulating.c88yhjcbemef.eu-north-1.rds.amazonaws.com:5432/formulating" -verbose up

echo "start app"
exec "$@"
