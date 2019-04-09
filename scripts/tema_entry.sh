#!/bin/sh

# Start hook in the background
/mws/temasync &

# Start the original entry point
/usr/local/bin/docker-entrypoint.sh "$@"