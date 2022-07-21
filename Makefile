postgres:
	docker run --name postgres12 --network formulating -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=eloh -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root formulating

dropdb:
	docker exec -it postgres12 dropdb formulating

migrateup:
	migrate -path db/migration -database "$DB_SOURCE" -verbose up

migrateup1:
	migrate -path db/migration -database "$DB_SOURCE" -verbose up 1

migratedown:
	migrate -path db/migration -database "$DB_SOURCE" -verbose down

migratedown1:
	migrate -path db/migration -database "$DB_SOURCE" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

generate:
	go generate ./...


.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server generate