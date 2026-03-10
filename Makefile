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
	