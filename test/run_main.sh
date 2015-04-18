#!/usr/bin/env bash

echo "Running main.go"; echo; echo;
while true
do
    go run main.go
    echo; echo; echo;
    echo "---------------------------- PRESS ENTER TO 'go run main.go' ----------------------------";
    read
    echo "Running main.go"; echo; echo;
done

