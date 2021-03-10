#!/bin/bash

docker run -it --rm -d --name binancebot -v $PWD/src:/go/src/binancebot golang:1.16 sh
