#!/bin/bash

echo "Checking and installing the uuid package..."
go get "code.google.com/p/go-uuid/uuid"

mkdir -p logs

echo "Installing module: client-app"
go install ./client-app

echo "Installing module: server-app"
go install ./server-app

echo "Installing module: space-app"
go install ./space-app
