migrateup:
	migrate -path db/migrations -database "postgresql://root:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test github.com/definitely-unique-username/simple_bank/db/sqlc -run TestMain

.PHONY: migrateup migratedown sqlc test
