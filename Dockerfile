FROM golang:alpine as builder

RUN apk update && apk add --no-cache git

COPY . $GOPATH/src/github.com/gosu-team/fptu-api/
WORKDIR $GOPATH/src/github.com/gosu-team/fptu-api/

RUN go get -v -d ./cmd/app/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/app/main.go

CMD ["./main"]

# This is docker build command: 
# sudo docker build -t fptu-api .

# This is docker run command:
# sudo docker run -dit -p 3000:3000 fptu-api:latest