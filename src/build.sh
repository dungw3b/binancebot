#!/bin/bash

env GOOS=linux GOARCH=amd64 go build -v -o ../bin/bbot-linux-amd64 main.go binance.go util.go
env GOOS=darwin GOARCH=amd64 go build -v -o ../bin/bbot-darwin-amd64 main.go binance.go util.go
env GOOS=windows GOARCH=amd64 go build -v -o ../bin/bbot-windows-amd64 main.go binance.go util.go