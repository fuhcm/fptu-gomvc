#!/bin/sh
ssh root@gosu.team <<EOF
    cd /home/golang/src/github.com/gosu-team/fptu-api
    git checkout .
    git pull
    docker build -t fptu-api .
    docker stop fptu-api
    docker rm fptu-api
    docker run -d --name fptu-api -p 5001:3000 fptu-api:latest
    exit
EOF
