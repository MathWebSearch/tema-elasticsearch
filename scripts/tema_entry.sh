#!/bin/bash

# Start temasearch
/usr/local/bin/docker-entrypoint.sh "$@" &
ELASTICPID=$!

/mws/tema-elasticsync &
SYNCPID=$!

# Trap signals to kill elastic sync
trap 'kill -TERM $SYNCPID' TERM INT

# wait for the elasticsync to exit and reset the trap
wait $SYNCPID
trap 'kill -TERM $ELASTICPID' TERM INT
wait $SYNCPID
SYNCSTATUS=$?


# Syncronize and exit if it failed
if [ $SYNCSTATUS -eq 0 ];
then
    echo "Sync success, handing control to elasticsearch"
else
    echo "Sync failed, stopping elasticsearch"
    kill -TERM $ELASTICPID
fi

# Wait for the pid etc
wait $ELASTICPID
trap - TERM INT
wait $ELASTICPID
exit $?