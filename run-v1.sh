#!/bin/bash
URL = $1
git checkout master
echo "Migrate all down to lowest version . . ."
migrate -path db/migration -database $1 -verbose down
echo "Migrate all down to v1 . . ."
migrate -path db/migration -database $1 -verbose up 1
echo "Inserting mock data"
go run insertMock.go
echo "Running server . . ."
go run main.go