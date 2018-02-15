#!/bin/bash

### PARAMS
file="1userWorkLoad"
rate="100" # in ms

# Little animation before we execute
for i in `seq 5 -1 1`; do 
    echo -en "\rExecuting $file in $i seconds"
    sleep 1
done 
echo -e "\nStarting\n"

go run workgen.go -f files/$file -r $rate