#!/bin/bash
URL = $1
git checkout v2
echo "Migrate up to v3 . . ."
migrate -path db/migration -database $1 goto 2
migrate -path db/migration -database $1 -verbose up 1
echo "Running server . . ."
go run main.go