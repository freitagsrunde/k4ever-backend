VERSION=0.0.2

all: build

build:
	go build -ldflags "-X github.com/freitagsrunde/k4ever-backend/internal/context.GitCommit=$$(git rev-parse HEAD) -X github.com/freitagsrunde/k4ever-backend/internal/context.GitBranch=$$(git rev-parse --abbrev-ref HEAD) -X github.com/freitagsrunde/k4ever-backend/internal/context.BuildTime=$$(date -u --iso-8601=seconds) -X github.com/freitagsrunde/k4ever-backend/internal/context.version=${VERSION}"

run:
	go run -ldflags "-X github.com/freitagsrunde/k4ever-backend/internal/context.GitCommit=$$(git rev-parse HEAD) -X github.com/freitagsrunde/k4ever-backend/internal/context.GitBranch=$$(git rev-parse --abbrev-ref HEAD) -X github.com/freitagsrunde/k4ever-backend/internal/context.BuildTime=$$(date -u --iso-8601=seconds) -X github.com/freitagsrunde/k4ever-backend/internal/context.version=${VERSION}"

docker-build:
	docker build -t freitagsrunde/k4ever-backend -t freitagsrunde/k4ever-backend:${VERSION} .

deploy:
	scripts/docker_push.sh ${VERSION}
