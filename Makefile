.PHONY: up down build test run migrate proto

up:
	docker-compose up -d --build

down:
	docker-compose down

build:
	go build -o bin/auth-service auth-service/main.go
	go build -o bin/task-service task-service/main.go
    go build -o bin/user-service user-service/main.go
	go build -o bin/notification-service notification-service/main.go
    go build -o bin/gateway-service gateway-service/main.go


test:
	go test -v ./...

run:
    go run gateway-service/main.go

migrate:
	docker exec -it your-postgres-container /bin/bash -c "psql -U your-db-user -d your-db-name -f /docker-entrypoint-initdb.d/init.sql"

proto:
    protoc --go_out=. --go-grpc_out=. auth-service/pb/auth.proto
    protoc --go_out=. --go-grpc_out=. task-service/pb/task.proto
    protoc --go_out=. --go-grpc_out=. user-service/pb/user.proto
    protoc --go_out=. --go-grpc_out=. notification-service/pb/notification.proto
    protoc --go_out=. --go-grpc_out=. gateway-service/pb/*.proto

help:
	@echo "make up          - Run docker compose up"
	@echo "make down        - Run docker compose down"
	@echo "make build       - Build all the services"
	@echo "make test        - Run all the unit tests"
    @echo "make run         - Run gateway service"
	@echo "make migrate     - Run database migrations"
	@echo "make proto        - Generate proto files"