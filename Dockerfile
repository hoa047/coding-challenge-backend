FROM golang:1.9

RUN mkdir -p /go/src/src
RUN mkdir -p /go/src/main
WORKDIR /go/src/main

ADD src /go/src/src 
ADD wait-for-it.sh /go/src/main
ADD main/main.go /go/src/main

RUN go get -v