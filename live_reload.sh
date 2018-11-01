#!/bin/sh
# This script require 2 executable commands passed in as arguments
# arg1: main app which is reloaded when there are changes
# arg2: execute when change events are detected
# Example:
# ./live_reload.sh "bin/server" "make build-server"

sigint_handler()
{
  kill $PID
  exit
}

trap sigint_handler SIGINT

watch_files()
{
  # only watch for go files
  # suppress output
  fswatch -e ".*" -i "\\.go$" --one-event . > /dev/null
}

echo "#### Customed Live Reload!"
echo "#### When *.go files change, it will execute \"$2\"."
echo "#### Now executing \"$1\" ..."

while true; do
  $1 &
  PID=$!
  watch_files
  kill $PID
  # wait for PID to be shut down
  wait $PID 2>/dev/null

  # watch till $2 succeeds
  until $2; do
    watch_files
  done

  echo "Restarting $1 ..."
done
