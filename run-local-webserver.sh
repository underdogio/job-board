#!/usr/bin/env bash

source ./.env.dev

go build -o bin/server cmd/server/main.go

./bin/server
