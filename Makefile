BIN := "./bin/imgpreviewer"
DOCKER_IMG="imgpreviewer"
GIT_HASH := $(shell git log --format="%h" -n 1)

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.60.3
	golangci-lint --version

lint: install-lint-deps
	golangci-lint run ./...

test:
	go test -race -count 100 ./internal/...

integration-test:
	docker compose --file docker-compose.test.yaml up --detach --build
	go test -tags=integration ./tests/integration/
	docker compose --file docker-compose.test.yaml down

build:
	go build -v -o $(BIN) ./cmd/imgpreviewer

run:
	docker compose up --detach --build

build-image:
	docker build -t $(DOCKER_IMG):$(GIT_HASH) .

run-image: build-image
	docker run $(DOCKER_IMG):$(GIT_HASH)

.PHONY: build run build-image run-image test integration-test lint
