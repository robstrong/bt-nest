#!/bin/bash

GOOS=linux GOARCH=arm GOARM=6 go build -o bin/nest-bt-linux-arm6 .
GOOS=linux GOARCH=amd64 go build -o bin/nest-bt-linux-amd64 .
