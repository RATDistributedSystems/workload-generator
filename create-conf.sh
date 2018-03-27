#!/bin/bash

str=""
for i in $(echo $PARAM_WEBSERVER_ADDR | tr "," "\n"); do
    str="server $i:$PARAM_WEBSERVER_PORT;\n\t\t$str"
done

FILENAME=$1
sed -i "s/#PARAM_SERVER_CONFIG#/${str}/g" $FILENAME

/etc/init.d/nginx restart