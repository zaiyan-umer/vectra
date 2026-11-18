.PHONY: build run test test-race testserver testserver-build

build:
	@go -C app build -o bin/api ./cmd/loadbalancer

run: build
	@cd app && ./bin/api

test:
	@go -C app test ./...

test-race:
	@go -C app test -race ./...

testserver:
	@docker compose -f test-server/docker-compose.yml up -d

testserver-build:
	@docker compose -f test-server/docker-compose.yml up -d --build

migration-build:
	@go -C app build -o bin/migrate ./cmd/migrate

migrate-up: migration-build
	@cd app && ./bin/migrate up

migrate-down: migration-build
	@cd app && ./bin/migrate down