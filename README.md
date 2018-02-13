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
./workgen -f 1userWorkLoad
``` 

## Use with Docker

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

`./workgen -f 1userWorkLoad --ip 192.168.1.1`

The target port will default to `44440`. To modify use the `--port` flag

`./workgen -f 1userWorkLoad --port 44442`

The rate of execution can be specific with the `-r` flag 
which will add number of millisecond delay between commands. Default is 50 ms

`./workgen -f 1userWorkLoad -r 50`


Instead of point to a file, a single value can be passed in from the command line
with the `-c` flag

`./workgen -c "ADD, asinha94, 100.10"`

The workload generator can also be used as a generic TCP client with the `--TCP` flag.

`./workgen -f 1userWorkLoad --tcp=true` 

## Shorthand Options

The shorthand of every option is available e.g `--port` and `-port` are synonymous.