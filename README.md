# Go RESTful API for FPTU Tech Insider

This repository is a skeleton for building production ready Golang REST API.

In this repo, we used `gorilla/mux` with also pure Go.

## Development

Run development container:

`docker-compose up`

## Production

Build production image:

`docker build -t fptu-api .`

Run production container:

`docker run -d --name fptu-api -p 5001:3000 fptu-api:latest`

You may found the example environment file at `.env.example`.
