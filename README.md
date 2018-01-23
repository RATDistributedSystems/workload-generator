# Workload Generator

Generates HTTP POST requests directly to the results page of the webserver, which is where the transactions get handled

## How to use

build with rergular go tools

```
go get github.com/RATDistributedSystems
go build workgen.go
```

Execute by pointing to the worklaod file
```
./workgen 1000users.txt
``` 
## Optional command-line flags

By default the address to target will be localhost to modify, use the `--ip` flag

`./workgen 100users.txt --ip 192.168.1.1`

The port will default to `44440`. To modify use the `--port` flag

`./workgen 1user.txt --port 44442`

Both the `--ip` and `--port` can be used simultaneously. Both `-ip` and `-port` are also valid 