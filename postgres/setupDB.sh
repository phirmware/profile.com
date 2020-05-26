#! /usr/bin/env bash

docker rmi profile

docker build -f dockerfile -t profile .

docker run --rm -d  -p 5432:5432 -e POSTGRES_DB=profile_dev webdevgo:latest

echo "Run GO file main.go TO CHECK IF DB WAS SUCCESSFULLY CREATED"