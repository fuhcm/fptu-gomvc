# Go REST API for FPTU Tech Insider

This repository is a skeleton for building production ready Golang REST API.

In this repo, we used gorilla/mux.

# Installation

Assuming you have a working Go environment and GOPATH/bin is in your PATH.

## Dependencies

No need to fetch dependencies, thanks to `govendor`.

Generate & modify your own environment configuration file:

```
$ mv .env.example .env
$ vim .env
```

Start server proxy to listen on port 3000 and send request to proxied app listening on port 8000:

```
$ go run cmd/app/main.go
```
