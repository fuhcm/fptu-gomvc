# Go REST API for FPTu.tech

This repository is a skeleton for building production ready Golang REST API.

In this repo, we used `gorilla/mux` with almost pure Go.

## Development

Run development container:

`$ docker-compose up`

## Production

Build production image:

`$ docker build -t fptu-api .`

Run production container:

`$ docker run -d --name fptu-api -p 5001:3000 fptu-api:latest`

You may found the example environment configuration file at `.env.example`.
