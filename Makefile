include .env
export

migrate-up:
	migrate -database $(CONN_STRING) -path migrations up
run:
	go run main.go