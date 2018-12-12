#!/bin/bash
rm -rf main
go build cmd/app/main.go
fuser -n tcp -k 5001
nohup ./main &