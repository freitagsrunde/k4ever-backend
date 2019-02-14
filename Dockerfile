FROM golang AS builder

WORKDIR $GOPATH/src/github.com/freitagsrunde/k4ever-backend/
COPY . .

RUN curl -L -s https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 -o /go/bin/dep && \
    chmod +x /go/bin/dep && \
    dep ensure && \
    CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w -extldflags -static' -o /go/bin/k4ever

FROM alpine

COPY --from=builder /go/bin/k4ever /go/bin/k4ever

RUN chmod +x /go/bin/k4ever

ENTRYPOINT ["/go/bin/k4ever"]



