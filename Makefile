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

.PHONY: swag_gen
swag_gen:
	@echo "Generating swagger docs..."
	swag init -d cmd/server -o docs

.PHONY: swagger
swagger:
	@echo "Starting swagger docs..."
	docker compose up -d swagger