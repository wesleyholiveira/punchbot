FROM golang:alpine

WORKDIR /go/src/github.com/wesleyholiveira/punchbot
COPY . .

RUN apk add git && \
    go get -u github.com/kardianos/govendor && \
    govendor add +local +external +vendor && \
    govendor fetch +missing && \
    go build -o /go/bin/punchbot

ENTRYPOINT ["/go/bin/punchbot"]
