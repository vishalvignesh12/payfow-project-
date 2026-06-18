.PHONY: run stop build test migrate seed

run:
	docker-compose up --build

stop:
	docker-compose down

build:
	go build -o bin/server ./cmd/server

test:
	go test ./... -v

migrate:
	docker-compose exec api go run ./cmd/migrate

seed:
	docker-compose exec api go run ./cmd/seed