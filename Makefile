DB_V ?= 1

.PHONY: db-up
db-up:
	docker run --name=balance_service -e POSTGRES_PASSWORD=qwerty -p 5432:5432 -d --rm postgres

.PHONY: mgt
mgt:
	migrate -path ./migrations -database "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable" up
	
.PHONY: mgt-frc
mgt-frc:
	@echo "migrage force для БД Версии DB_V: $(DB_V)"
	migrate -path ./migrations -database "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable" force $(DB_V)

.PHONY: gen-pb
gen-pb:
	@echo -e "=== protoc: Компиляция balanceSrvc.proto ===\n" 

	protoc -I ./protos/proto ./protos/proto/balanceSrvc.proto \
	--go_out=./protos/gen --go_opt=paths=source_relative \
	--go-grpc_out=protos/gen --go-grpc_opt=paths=source_relative
	
	@echo -e "\n=== protoc: Компиляция завершена ==="

.PHONY: cover
cover:
	go test -v -short -count=1 -race -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: integration-test
integration-test:
	@echo "Запуск окружения..."
	@docker run --name=balance_service-test -e POSTGRES_PASSWORD=qwerty -p 5433:5432 -d --rm postgres;
	@trap 'echo "Остановка контейнера..."; docker rm -f balance_service-test' EXIT; \
	echo "Ожидание готовности БД..."; \
	sleep 5; \
	echo "Запуск миграций..."; \
	migrate -path ./migrations -database "postgres://postgres:qwerty@127.0.0.1:5433/postgres?sslmode=disable" up;\
	echo "Запуск тестов..."; \
	go test -v ./tests/

.PHONY: run-app
run-app:
	make db-up; 
	@trap 'echo "Остановка контейнера..."; docker rm -f balance_service' EXIT; \
	sleep 5; \
	make mgt; \
	go run cmd/balanceSrvc/main.go -config ./config/local.yml
