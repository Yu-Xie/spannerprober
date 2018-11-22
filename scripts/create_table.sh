#!/usr/bin/env bash

GOOGLE_APPLICATION_CREDENTIALS=config/credentials.json \
    SOURCE_REGION=$1 \
    SPANNER_NAME=$2 \
    CREATE=true \
    go run main.go
