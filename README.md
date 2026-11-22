# Vectra

A Go HTTP load balancer with 5 routing algorithms, health checking, autocannon benchmarking, and a React dashboard.

## Features

- **5 routing algorithms**: round_robin, healthy_round_robin, least_connections, latency_aware, adaptive_selection
- **Health checks** with configurable intervals
- **Autocannon benchmarking** with per-backend traffic distribution
- **Real-time dashboard** (React + TypeScript + Vite) showing RPS, latency, error rate, backend allocation
- **PostgreSQL** for benchmark history

## Quick Start

```bash
# 1. Start test backends (3 servers on 8081-8083)
make testserver

# 2. Configure environment
cp app/.env.example app/.env
# Edit app/.env with your DATABASE_URL

# 3. Run database migrations
make migrate-up

# 4. Start load balancer on :8080
make run

# 5. Start dashboard on :5173
cd dashboard && cp .env.example .env && bun install && bun run dev
```

## Algorithms

- `round_robin` — Sequential distribution across all backends
- `healthy_round_robin` — Round robin skipping unhealthy backends
- `least_connections` — Route to backend with fewest active connections
- `latency_aware` — Prefer backends with lower latency
- `adaptive_selection` — Combines health, latency, and connection count

## Benchmarking

Open the dashboard → "New Benchmark" → set duration (10-120s) and connections (1-1000). Results show request rate, latency percentiles, status codes, and per-backend traffic distribution.

## Configuration

See `app/.env.example` and `dashboard/.env.example` for all options.