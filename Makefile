include .env

migrateup:
	@goose -dir internal/sql/schema postgres "postgresql://${DB_DOCKER_USER}:${DB_DOCKER_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" up

migratedown:
	@goose -dir internal/sql/schema postgres "postgresql://${DB_DOCKER_USER}:${DB_DOCKER_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" down
	
run:
	go build -o bin/api cmd/api/main.go && bin/api
	
test:
	@go test -v ./...

create_db:
	docker-compose exec db createdb --username=${DB_DOCKER_USER} --owner=${DB_DOCKER_USER} ${DB_NAME}