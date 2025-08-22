.PHONY: up
up:
	@docker compose up -d --build

.PHONY: down
down:
	@docker compose down

.PHONY: clean
clean: down
	@docker volume rm -f goexpert-rate-limiter_redis_data 2>/dev/null || true
