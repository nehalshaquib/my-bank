postgres:
	docker run --name my-bank-postgres15 -p 5434:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it my-bank-postgres15 createdb --username=root --owner=root my_bank

dropdb:
	docker exec -it my-bank-postgres15 dropdb my_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5434/my_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5434/my_bank?sslmode=disable" -verbose down

.PHONY: createdb dropdb migrateup migratedown postgres