# board
#
# VERSION               1.0
FROM golang:alpine

RUN apk update
RUN apk add git

ADD . /go/src/board
WORKDIR /go/src/board

RUN go get -v ./...

EXPOSE 8080
ENTRYPOINT board -p 8080
