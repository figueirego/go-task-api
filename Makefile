db-up:
	docker compose up -d

db-down:
	docker compose down

db-reset:
	docker compose down -v
	docker compose up -d

migrate:
	go run ./cmd/migrate

api:
	go run ./cmd/api

test:
	go test ./...

test-race:
	go test -race ./...

fmt:
	go fmt ./...