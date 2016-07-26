FROM golang:1.6
MAINTAINER meomap

ENV GOOS=windows

ADD . /go/src/github.com/meomap/logruswindows

RUN go get -v github.com/meomap/logruswindows
