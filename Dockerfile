# board
#
# VERSION               1.0
FROM golang:alpine

ADD . /go/src/board
WORKDIR /go/src/board

RUN go build .

EXPOSE 8080
ENTRYPOINT ./board -p 8080
