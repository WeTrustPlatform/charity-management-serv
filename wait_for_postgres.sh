#!/bin/bash
# wait-for-postgres.sh
# Courtesy of https://docs.docker.com/compose/startup-order/
# https://stackoverflow.com/questions/4922943/test-from-shell-script-if-remote-tcp-port-is-open

set -e
cmd="$@"
source .env

poll()
{
  timeout 1 bash -c "cat < /dev/null > /dev/tcp/$DB_HOST/$DB_PORT"
}

postgres="Postgres $DB_HOST:$DB_PORT"
until poll; do
  >&2 echo "$postgres is unavailable - sleeping"
  sleep 2
done

echo "$postgres is up - executing command $cmd"
$cmd
