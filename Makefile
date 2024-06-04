run:
	go run cmd/server/main.go

start:
	@docker-compose up --build -d
	@echo Starting Go application..
	@go run cmd/server/main.go

stop:
	@docker compose down
	@echo All done!

test:
	@go test -v ./...

test-coverage:
	@go test -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out
