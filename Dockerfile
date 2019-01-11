FROM golang:alpine as builder

RUN apk update && apk add --no-cache git

COPY . $GOPATH/src/github.com/gosu-team/fptu-api/
WORKDIR $GOPATH/src/github.com/gosu-team/fptu-api/

RUN go get -v -d ./cmd/app/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/app/main.go

FROM ellerbrock/alpine-curl-ssl

WORKDIR /root/src/app

COPY --from=builder /go/src/github.com/gosu-team/fptu-api/main /root/src/app/main
COPY --from=builder /go/src/github.com/gosu-team/fptu-api/.env /root/src/app/.env

EXPOSE 3000

ENTRYPOINT ["./main"]

# This is docker build command: 
# docker build -t fptu-api .

# This is docker run command:
# docker run -d --name fptu-api -p 5001:3000 fptu-api:latest