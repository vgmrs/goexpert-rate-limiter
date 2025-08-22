# goexpert-rate-limiter

## Description

A rate limiter middleware for Go applications that can limit requests based on IP address or API tokens.

## How to Install

1. Make sure you have Go installed (version 1.20 or higher)

2. Start containers:
   ```bash
   make up
   ```

## How to Run

1. The rate limiter will be available at `http://localhost:8080`

3. Test the rate limiter by making requests to any endpoint. The limiter can be configured to work with:
   - IP-based limiting
   - API token-based limiting (use header: `API_KEY: your-token`)

## Configuration

Configure the rate limiter using environment variables at `docker-compose.yml` file:

```
RATE_LIMIT_IP=10
RATE_LIMIT_TOKEN=100
BLOCK_DURATION_SECONDS=300
REDIS_ADDR=redis:6379
```

- `RATE_LIMIT_IP`: Maximum requests per second per IP
- `RATE_LIMIT_TOKEN`: Maximum requests per second per token
- `BLOCK_DURATION_SECONDS`: Block duration after rate limit is exceeded
- `REDIS_ADDR`: Redis server address

## Other available Commands

- `make down`: Stop all services
- `make clean`: Stop services and remove Redis volume
