#!/bin/bash
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push freitagsrunde/k4ever-backend

if [[ -n $1 ]]; then
    docker push freitagsrunde/k4ever-backend:$1
fi
