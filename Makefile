.PHONY: up down logs migrate-up migrate-down migrate-down-all

DB_DSN = postgres://postgres:postgres@postgres:5432/eff_db?sslmode=disable

up:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f api

migrate-up:
	docker compose --profile migrate run --rm migrate \
		-path=/migrations -database=$(DB_DSN) up

migrate-down:
	docker compose --profile migrate run --rm migrate \
		-path=/migrations -database=$(DB_DSN) down 1

migrate-down-all:
	docker compose --profile migrate run --rm migrate \
		-path=/migrations -database=$(DB_DSN) down -all
