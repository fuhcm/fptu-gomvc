FROM golang:alpine as builder

COPY . $GOPATH/src/webserver/
WORKDIR $GOPATH/src/webserver/

RUN touch .env

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

FROM alpine

RUN apk --no-cache add ca-certificates rsync openssh

WORKDIR /root/src/app

COPY --from=builder /go/src/webserver/main /root/src/app/main
COPY --from=builder /go/src/webserver/.env /root/src/app/.env

EXPOSE 3000

ENTRYPOINT ["./main"]

# This is docker build command: 
# docker build -t fptu-api .

# This is docker run command:
# docker run -d --name fptu-api -p 5001:3000 fptu-api:latest