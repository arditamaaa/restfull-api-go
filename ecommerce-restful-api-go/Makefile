include .env
export $(shell sed 's/=.*//' .env)

start:
	@go run main.go

create_migration:
	@migrate create -ext=sql -dir=database/migrations -seq create_$(name)_table

migrate_up:
	migrate -path=database/migrations -database "mysql://${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?charset=utf8mb4&parseTime=True&loc=Local" -verbose up $(version)

migrate_down:
	migrate -path=database/migrations -database "mysql://${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?charset=utf8mb4&parseTime=True&loc=Local" -verbose down $(version)

migrate_up_seed:
	migrate -path=database/migrations -database "mysql://${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?charset=utf8mb4&parseTime=True&loc=Local" -verbose up $(version) && go run database/seeder/main.go

seed:
	go run database/seeder/main.go
