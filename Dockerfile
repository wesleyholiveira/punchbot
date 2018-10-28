FROM golang:alpine

WORKDIR /go/src/github.com/wesleyholiveira/punchbot
COPY . .

RUN apk add git && \
    go get -u github.com/kardianos/govendor && \
    go build -o /go/bin/punchbot

ENTRYPOINT ["/go/bin/punchbot"]
