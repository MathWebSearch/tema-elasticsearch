#!/bin/sh

# Start hook in the background
/mws/tema_hook &

# Start the original entry point
/usr/local/bin/docker-entrypoint.sh "$@"