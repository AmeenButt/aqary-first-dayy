migration_up: 
	migrate -path database/migration/ -database "postgresql://postgres:admin@localhost:5432/sqlc-practice?sslmode=disable" up
migration_down: 
	migrate -path database/migration/ -database "postgresql://postgres:admin@localhost:5432/sqlc-practice?sslmode=disable" down
run: 
	air
run-docker: 
	docker-compose up --build