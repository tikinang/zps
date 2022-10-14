#!/bin/bash

while true
do
  printf 'HTTP/1.1 200 OK\n\n%s' "$(cat index.html)" | nc -6 -v -l -N 1999
done
