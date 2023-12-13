include .env

create_migration:
	migrate create -ext sql -dir migrations create_$(table)

migration_up:
	migrate -path migrations -database "${DATABASE}://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" up

migration_rollback:
	migrate -path migrations -database "${DATABASE}://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" down