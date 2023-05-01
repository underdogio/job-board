#!/usr/bin/env bash
set -a; source .env; set +a

go build -o bin/server cmd/server/main.go

./bin/server
