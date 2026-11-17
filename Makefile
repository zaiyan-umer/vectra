.PHONY: build run test test-race testserver testserver-build

build:
	@go -C app build -o bin/api ./cmd/loadbalancer

run: build
	@./app/bin/api

test:
	@go -C app test ./...

test-race:
	@go -C app test -race ./...

testserver:
	@docker compose -f test-server/docker-compose.yml up -d

testserver-build:
	@docker compose -f test-server/docker-compose.yml up -d --build
