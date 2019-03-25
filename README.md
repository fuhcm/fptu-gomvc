# Go REST API for FPTU Tech Insider

This repository is a skeleton for building production ready Golang REST API.

In this repo, we used gorilla/mux.

## Development

Build develop image:

`docker build -f dev.Dockerfile -t fptu-api-dev .`

Run production container:

`docker run -it -p 3000:3000 --env-file ./.env -v $(pwd):/go/src/github.com/gosu-team/fptu-api fptu-api-dev:latest`

## Production

Build production image:

`docker build -t fptu-api .`

Run production container:

`docker run -d --name fptu-api -p 5001:3000 fptu-api:latest`
