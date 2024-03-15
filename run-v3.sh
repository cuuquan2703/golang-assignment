#!/bin/bash
URL = $1
git checkout v3
echo "Migrate all down to lowest version . . ."
migrate -path db/migration -database $1 -verbose down
echo "Migrate all down to v1 . . ."
migrate -path db/migration -database $1 -verbose up 3
echo "Running server . . ."
go run main.go