FROM golang:1.10.3 as builder

COPY . $GOPATH/src/github.com/gosu-team/fptu-api/
WORKDIR $GOPATH/src/github.com/gosu-team/fptu-api/

RUN go get -v -d ./cmd/app/

RUN env GOOS=linux GOARCH=amd64 go build cmd/app/main.go

FROM scratch

COPY --from=builder /go/src/github.com/gosu-team/fptu-api/main /go/bin/main

EXPOSE 3000

ENTRYPOINT ["/go/bin/main"]

# This is docker build command: 
# sudo docker build -t fptu-api .

# This is docker run command:
# sudo docker run -dit -p 3000:3000 fptu-api:latest