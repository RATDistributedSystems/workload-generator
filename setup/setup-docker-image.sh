#!/bin/bash

cd ..

# Build our app with dynamic linking disabled so it can be run with no libc support
CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo --ldflags="-s" -o workgen

# Build the image
docker build -t workgen .

# Delete the garbage
rm -rf workgen
