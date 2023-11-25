.PHONY: infra
infra:
	@echo "Starting up infrastructure..."
	docker compose up -d mongo auth-redis did-redis

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

.PHONY: gen
gen:
	@echo "Generating go files from abis..."
	mkdir -p pkg/contracts/gen/ticketex
	abigen --abi=pkg/contracts/abi/TicketExtended.abi --pkg=ticketex --out=pkg/contracts/gen/ticketex/ticketex.go
	mkdir -p pkg/contracts/gen/ticket
	abigen --abi=pkg/contracts/abi/Ticket.abi --pkg=ticket --out=pkg/contracts/gen/ticket/ticket.go
	mkdir -p pkg/contracts/gen/token
	abigen --abi=pkg/contracts/abi/HeroToken.abi --pkg=token --out=pkg/contracts/gen/token/token.go

.PHONY: swag_gen
swag_gen:
	@echo "Generating swagger docs..."
	swag init -d cmd/server -o docs --parseInternal --pdl 2

.PHONY: swagger
swagger:
	@echo "Starting swagger docs..."
	docker compose up -d swagger --build
	@echo "Swagger docs available at http://localhost:1323/swagger/index.html"