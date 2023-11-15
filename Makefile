.PHONY: up
up:
	@echo "Starting up containers..."
	docker compose up -d

.PHONY: up_build
up_build:
	@echo "Starting up containers..."
	docker compose up -d --build

.PHONY: down
down:
	@echo "Stopping containers..."
	docker compose down
