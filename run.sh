#!/bin/bash

nc -6 -l 2001 &

while true
do
  echo "running"
  sleep 10s
done
