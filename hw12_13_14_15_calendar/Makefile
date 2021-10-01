BIN := "./bin/calendar"
SCHEDULER := "./bin/calendar_scheduler"
SENDER := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

scheduler:
	go build -v -o $(SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	$(SCHEDULER)

sender:
	go build -v -o $(SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender
	$(SENDER)

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
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.37.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

enterdb:
	docker exec -it postgres psql -U homestead;

migrate: 
	goose -dir migrations postgres "user=homestead password=secret dbname=homestead sslmode=disable" up

rollback: 
	goose -dir migrations postgres "user=homestead password=secret dbname=homestead sslmode=disable" down

generate:
	protoc -I=api --go_out=internal/server/grpc/generated --go-grpc_out=internal/server/grpc/generated api/EventService.proto

evans:
	evans api/EventService.proto -p 2345
