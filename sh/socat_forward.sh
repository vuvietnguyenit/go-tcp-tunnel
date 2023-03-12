#!/bin/sh
source=$1
dest=$2
desc=$3
# Example: socat TCP4-LISTEN:4000,fork,reuseaddr TCP4:10.51.78.127:4000 &
socat TCP4-LISTEN:$source,fork,reuseaddr TCP4:$dest &
done