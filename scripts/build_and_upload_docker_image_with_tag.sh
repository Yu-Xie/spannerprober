#!/usr/bin/env bash

if [[ -z "$1" ]]
then
    echo "must set the tag"
    exit 1
fi

docker build -t spannerprober .

docker tag spannerprober gcr.io/horizon-spanner-benchmark/spannerprober:$1

docker push gcr.io/horizon-spanner-benchmark/spannerprober:$1
