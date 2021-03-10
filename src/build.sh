#!/bin/bash

env GOOS=linux GOARCH=amd64 go build -v -o ../bin/binancebot-linux-amd64 main.go