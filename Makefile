
all: build

build:
	go build -ldflags "-X github.com/freitagsrunde/k4ever-backend/internal/context.GitCommit=$$(git rev-parse HEAD) -X github.com/freitagsrunde/k4ever-backend/internal/context.GitBranch=$$(git rev-parse --abbrev-ref HEAD) -X github.com/freitagsrunde/k4ever-backend/internal/context.BuildTime=$$(date -u '+%Y-%m-%d_%I:%M:%S%p') -X github.com/freitagsrunde/k4ever-backend/internal/context.version=0.0.1"

run:
	go run -ldflags "-X github.com/freitagsrunde/k4ever-backend/internal/context.GitCommit=$$(git rev-parse HEAD) -X github.com/freitagsrunde/k4ever-backend/internal/context.GitBranch=$$(git rev-parse --abbrev-ref HEAD) -X github.com/freitagsrunde/k4ever-backend/internal/context.BuildTime=$$(date -u '+%Y-%m-%d_%I:%M:%S%p') -X github.com/freitagsrunde/k4ever-backend/internal/context.version=0.0.1"
