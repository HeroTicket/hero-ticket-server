.PHONY: up
up:
	@echo "Starting up containers..."
	docker compose up -d

.PHONY: down
down:
	@echo "Stopping containers..."
	docker compose down
