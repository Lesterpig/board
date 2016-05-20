# board
#
# VERSION               1.0
FROM golang:alpine

RUN apk update
RUN apk add git
RUN go get github.com/miekg/dns
RUN go build github.com/miekg/dns

ADD . /go/src/board
WORKDIR /go/src/board

RUN go build .

EXPOSE 8080
ENTRYPOINT ./board -p 8080
