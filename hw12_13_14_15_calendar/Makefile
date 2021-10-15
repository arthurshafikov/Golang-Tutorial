BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test --short -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.37.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

enterdb:
	docker exec -it postgres psql -U homestead;

generate:
	protoc -I=api --go_out=internal/server/grpc/generated --go-grpc_out=internal/server/grpc/generated api/EventService.proto

up:
	docker-compose up --build

down:
	docker-compose down --remove-orphans --volumes

integration-tests:
	docker-compose -f docker-compose.tests.yml up --build --abort-on-container-exit --exit-code-from integration-tests
	docker-compose -f docker-compose.tests.yml down --volumes

reset-integration-tests:
	docker-compose -f docker-compose.tests.yml down --volumes	
