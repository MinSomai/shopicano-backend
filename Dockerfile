FROM golang:alpine AS builder

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN apk add git openssh

ENV GOPATH=/go

ENV GOOS="linux"
ENV GOARCH="amd64"
ENV GO111MODULE=on

COPY . $GOPATH/src/github.com/shopicano/shopicano-backend
WORKDIR $GOPATH/src/github.com/shopicano/shopicano-backend

RUN go get github.com/ugorji/go@v1.1.2-0.20180831062425-e253f1f20942

RUN go get .
RUN rm /go/pkg/mod/github.com/coreos/etcd@v3.3.10+incompatible/client/keys.generated.go
RUN cp ./hacks/keys.generated.go /go/pkg/mod/github.com/coreos/etcd@v3.3.10+incompatible/client/

RUN go build -v -o shopicano
RUN mv shopicano /go/bin/shopicano

FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /root

COPY --from=builder /go/bin/shopicano /usr/local/bin/shopicano

ENTRYPOINT ["shopicano"]
