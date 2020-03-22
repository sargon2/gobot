#!/bin/bash

# TODO tag it with a temporary, unique local version
docker build -t gobot --build-arg USER_ID=$(id -u) --build-arg GROUP_ID=$(id -g) .

docker run -v $PWD:/opt/mount --rm --entrypoint cp gobot /gobot/go.{mod,sum} /opt/mount/
