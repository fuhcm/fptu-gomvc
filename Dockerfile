# Run these command to create docker image & run docker container from image
# docker build -t goapi .
# sudo docker run -p 3000:3000 goapi

FROM golang:1.10.3

RUN mkdir -p usr/src/app

WORKDIR usr/src/app

COPY main .
COPY .env .

CMD ["./main"]

EXPOSE 3000