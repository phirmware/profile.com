#! /usr/bin/env bash

RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

image_tag="profile"
port=5432
db_name="profile_dev"

printOk() {
    echo -e "$BLUE $1 $NC"
}

printBad() {
    echo -e "$RED $1 $NC"
}

checkDB() {
    FILE=main.go
    if test -f "$FILE"; then
        printOk "$FILE exists"
    else
        printBad "$FILE not found"
    fi

    go run $FILE
    if [ $? -eq 0 ]; then
        printOk "DB is Successfully setup, Happy Coding"
    else
        printBad "Something went wrong, try running main.go file manually"
    fi
}

docker rmi $image_tag
if [ $? -eq 0 ]; then
    printOk "Image found, removing image from system"
else
    printBad "Image not found, creating new image"
fi

docker build -f dockerfile -t $image_tag .
if [ $? -eq 0 ]; then
    printOk "Image Built Successfully"
else
    printBad "Something went wrong, try again"
fi

docker run --rm -d -p $port:$port -e POSTGRES_DB=$db_name $image_tag:latest
if [ $? -eq 0 ]; then
    printOk "Container is running, checking for successful DB setup"
    sleep 5
    checkDB
else
    printBad "Couldnt start up the container, please try again"
fi
