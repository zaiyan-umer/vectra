# Vectra

Vectra is an HTTP reverse proxy that distributes requests across configured
backend servers using round-robin selection.

## Current status

The application currently supports:

- Round-robin backend selection
- HTTP and HTTPS backend URLs
- Reverse-proxy forwarding
- `503 Service Unavailable` when no selected backend exists
- A Dockerized test backend with health, delay, and error endpoints

Health-aware routing, configuration files, graceful shutdown, and adaptive
selection are planned next.

## Run locally

Start the load balancer from the `app` directory:

```sh
go run ./cmd/loadbalancer
```

It listens on `:8080` and forwards to `localhost:8081`, `localhost:8082`,
and `localhost:8083`. Start the test backends separately with:

```sh
cd test-server
docker compose up --build
```

## Test

```sh
cd app
go test ./...
go test -race ./...
```
