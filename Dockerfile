FROM golang:alpine as builder
RUN apk add ca-certificates git

RUN mkdir -p /root/src/go
WORKDIR /root/src/go

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

FROM alpine

RUN apk add ca-certificates rsync openssh

WORKDIR /root/src/go

COPY --from=builder /root/src/go/main /root/src/go/main

EXPOSE 5100

ENTRYPOINT ["./main"]

# This is docker build command: 
# docker build -t fptu-api .

# This is docker run command:
# docker run -d --name fptu-api -p 5001:3000 fptu-api:latest