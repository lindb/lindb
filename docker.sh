#!/bin/bash

all_before(){
    echo "### make build-all start!"
    make build GOOS=linux
    echo "### make build-all done!"
}

build(){
    all_before
    echo "docker build start"
    docker-compose build
    echo "docker build done"
}

start(){
    echo "docker start"
    docker-compose up
    echo "start docker done"
}

stop(){
    echo "docker stop"
    docker-compose stop
    echo "stop docker done"
}

help(){
    echo "help    Display this help."
    echo "build   Build project docker image."
    echo "start   start project docker image."
    echo "stop    start project docker image."
}

case "$1" in
    build)
        build
        ;;
    help)
        help
        ;;
    start)
        start
        ;;
    stop)
        stop
        ;;
    *)
        help
        exit 1
        ;;
esac
exit 0