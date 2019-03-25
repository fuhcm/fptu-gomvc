# Go REST API for FPTU Tech Insider

This repository is a skeleton for building production ready Golang REST API.

In this repo, we used gorilla/mux.

## Development

Build develop image:

`docker-compose build`

Run production container:

`docker-compose up`

## Production

Build production image:

`docker build -t fptu-api .`

Run production container:

`docker run -d --name fptu-api -p 5001:3000 fptu-api:latest`
