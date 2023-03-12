#!/bin/sh

port=$1
cmd=`ss -tlnp | awk '{print $4}' | grep :$port`
if [ -n "$cmd" ]; then
    return 1
else
    return 0
fi
