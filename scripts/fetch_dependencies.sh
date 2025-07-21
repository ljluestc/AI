#!/bin/bash

# Script to fetch all project dependencies

echo "Fetching Go dependencies..."
go mod tidy
go get github.com/go-git/go-git/v5
go get github.com/neo4j/neo4j-go-driver/v4/neo4j
go get github.com/gorilla/mux
go get github.com/stretchr/testify
go get gonum.org/v1/gonum

echo "Dependencies updated successfully"
