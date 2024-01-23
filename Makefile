migration_up: 
	migrate -path database/migration/ -database "postgresql://postgres:admin@localhost:5432/sqlc-practice?sslmode=disable" up
migration_down: 
	migrate -path database/migration/ -database "postgresql://postgres:admin@localhost:5432/sqlc-practice?sslmode=disable" down
run: 
	powershell -Command "(Get-Content .env) -replace 'DB_URL=.*', 'DB_URL=postgresql://postgres:admin@localhost:5432/sqlc-practice' | Set-Content .env"
	air
run-docker: 
	powershell -Command "(Get-Content .env) -replace 'DB_URL=.*', 'DB_URL=postgresql://postgres:admin@host.docker.internal:5432/sqlc-practice' | Set-Content .env"
	docker-compose up --build
sqlc:
	sqlc generate