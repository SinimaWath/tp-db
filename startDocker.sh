#!/usr/bin/env bash

make download-generators
docker build -t tp-db -f deploy/Dockerfile.golang .
./deploy/runDocker.sh