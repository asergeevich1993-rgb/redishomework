include .env
export

migrate-up:
	migrate -database "$(CONN_STRING)" -path migrations up
run:
	docker compose up -d --build