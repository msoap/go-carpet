#!/usr/bin/bash

args=""
while [ $# -gt 0 ]; do
    if [[ "$1" == -* ]]; then
        args+="$1"
    else
        space_files=$@
        files=${space_files// /,}
        break
    fi
    shift
done

go-carpet -summary -enforce $args -file $files
