#!/bin/bash

registry(){
    docker run -d --rm -p 5000:5000 --name registry registry:2
}

test(){
    ./dip -image sonatype/nexus3 -registry localhost:5000/ -latest "^(\d+\.){2}\d$" -preserve
    ./dip -image traefik -registry localhost:5000/ -latest "^v(\d+\.){1,2}\d+$" -preserve
}

cleanup(){
    docker stop registry
}

main(){
    registry
    test
}

trap cleanup EXIT
main