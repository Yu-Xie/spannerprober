#!/usr/bin/env bash

docker build -t spannerprober .

docker run --rm spannerprober -i -t --name spannerprober

