#!/bin/sh
ssh root@gosu.team <<EOF
    fuser -n tcp -k 5001
    cd /home/golang/src/github.com/gosu-team/fptu-api
    git checkout .
    git pull
    rm -rf main
    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=/home/golang
    go build cmd/app/main.go
    exit
EOF
cd /home/golang/src/github.com/gosu-team/fptu-api
nohup ./main > output.log&