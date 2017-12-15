#!/usr/bin/env bash

docker run --rm \
  -e COMPRESS_BINARY=true \
  -e LDFLAGS='-extldflags "-static"' \
  -v "$(pwd):/src" \
  -v /var/run/docker.sock:/var/run/docker.sock \
  centurylink/golang-builder \
  lmouline/calculator:0.1