postgres:
	docker run --name myBank-postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it myBank-postgres15 createdb --username=root --owner=root myBank

dropdb:
	docker exec -it myBank-postgres15 dropdb myBank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/myBank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/myBank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: createdb dropdb migrateup migratedown postgres server