FROM golang:alpine
RUN apk add ca-certificates git

RUN mkdir -p /root/src/go
WORKDIR /root/src/go
COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 5000

ENTRYPOINT ["go","run","main.go"]

# This is docker build command: 
# docker build -f dev.Dockerfile -t fptu-api-dev .

# This is docker run command:
# docker run -it -p 3000:3000 --env-file ./.env -v $(pwd):/go/src/github.com/gosu-team/fptu-api fptu-api-dev:latest