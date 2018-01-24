#!/bin/bash

cd ..

# Build our app in docker with CGO disabled (dynamic linking) then copy it back outside
docker run --rm -it -v "$GOPATH":/gopath -v "$(pwd)":/app -e "GOPATH=/gopath" -w /app golang:1.9 sh -c 'CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o workgen'

# Build the image
docker build -t workgen .

# Delete the garbage
rm -rf workgen
