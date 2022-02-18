#!/bin/bash

echo "Creating Adapter: $1";

cd ../
mkdir $1
cd $1
mkdir cmd
touch cmd/main.go
mkdir config
mkdir config/sample
mkdir pkg
mkdir pkg/adapter
touch pkg/adapter/adapter.go
touch Dockerfile
touch README.md
go mod init
go mod tidy
