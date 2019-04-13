FROM golang AS builder

WORKDIR $GOPATH/src/github.com/freitagsrunde/k4ever-backend/
COPY . .

RUN curl -L -s https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 -o /go/bin/dep && \
    go get -u github.com/go-swagger/go-swagger/cmd/swagger && \
    go get -u github.com/gobuffalo/packr/packr && \
    chmod +x /go/bin/dep && \
    dep ensure && \
    go generate && \
    packr && \
    CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags "-w -extldflags -static -X github.com/freitagsrunde/k4ever-backend/internal/context.GitCommit=$(git rev-parse HEAD) -X github.com/freitagsrunde/k4ever-backend/internal/context.GitBranch=$(git rev-parse --abbrev-ref HEAD) -X github.com/freitagsrunde/k4ever-backend/internal/context.BuildTime=$(date -u '+%Y-%m-%d_%I:%M:%S%p') -X github.com/freitagsrunde/k4ever-backend/internal/context.version=0.0.1" -o /go/bin/k4ever

FROM alpine:3.5

COPY --from=builder /go/bin/k4ever /go/bin/k4ever

RUN chmod +x /go/bin/k4ever && apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ENTRYPOINT ["/go/bin/k4ever"]



