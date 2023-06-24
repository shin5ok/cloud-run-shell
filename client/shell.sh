#!/bin/bash

AUTH_HEADER="X-MyGCP-Secret"
SECRET=gcp

if [ $# -lt 0 ] ;
then
    cat << EOM
    Usage:
    \$ $0 command 
    ex:
    \$ $0 ls -l
    \$ $0 ps axw
EOM
    exit 0
fi

CMD="{\"command\": \"$@\"}"
echo $CMD
echo $CMD | curl -s -H "Content-Type: application/json" -H "${AUTH_HEADER}: $SECRET" $URL/shellcommand -d @-
