# Workload Generator

Generates HTTP POST requests directly to the results page of the webserver, which is where the transactions get handled

## How to use

build with rergular go tools

```
go get github.com/RATDistributedSystems/workload-generator
cd $GOHOME/src/github.com/RATDistributedSystems/workload-generator
go build workgen.go
```

Execute by pointing to the worklaod file
```
./workgen 1000users.txt
``` 

## Use with Docker

The milestones are all executed in Docker, so we have generated a few docker images

1. Building image yourself

To build the docker image yourself, use the following commands

```
cd setup
./setup-docker-image.sh
```

This can then be executed as

```
docker run workgen
```

if you are using the same computer to run the webserver you will need to join networks

```
docker run --network container:<name> workgen
```
where `<name>` is the name or id of the webserver container

## Optional command-line flags

By default the target address will be localhost. To modify, use the `--ip` flag

`./workgen 100users.txt --ip 192.168.1.1`

The target port will default to `44440`. To modify use the `--port` flag

`./workgen 1user.txt --port 44442`

Both the `--ip` and `--port` can be used simultaneously. `-ip` and `-port` are also valid 