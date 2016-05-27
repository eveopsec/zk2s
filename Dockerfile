FROM golang:alpine

RUN apk add --no-cache git

RUN go get github.com/eveopsec/zk2s

WORKDIR /go/src/github.com/eveopsec/zk2s

RUN go get

RUN go build

ENTRYPOINT /go/src/github.com/eveopsec/zk2s/zk2s start
