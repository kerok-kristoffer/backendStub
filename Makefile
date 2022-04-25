postgres:
	docker run --name postgres12 -p 5454:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=eloh -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root formulating

dropdb:
	docker exec -it postgres12 dropdb formulating

migrateup:
	migrate -path db/migration -database "postgresql://root:eloh@localhost:5454/formulating?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:eloh@localhost:5454/formulating?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc