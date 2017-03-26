FROM golang:1.8.0-alpine

ADD . /go/src/github.com/Tinker-Ware/gh-service

RUN go install github.com/Tinker-Ware/gh-service   

ENTRYPOINT /go/bin/gh-service

EXPOSE 3000