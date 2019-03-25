FROM golang:alpine as builder

COPY . $GOPATH/src/github.com/gosu-team/fptu-api/
WORKDIR $GOPATH/src/github.com/gosu-team/fptu-api/

EXPOSE 3000

ENTRYPOINT ["go","run","cmd/app/main.go"]

# This is docker build command: 
# docker build -f dev.Dockerfile -t fptu-api-dev .

# This is docker run command:
# docker run -it -p 3000:3000 --env-file ./.env -v $(pwd):/go/src/github.com/gosu-team/fptu-api fptu-api-dev:latest