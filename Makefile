postgres:
	docker run --name=db -ti -e POSTGRES_PASSWORD='root' -p 5436:5432 -d --rm postgres:12.11

migrate:
	migrate -source file://migrations -database postgres://postgres:root@localhost:5436/postgres?sslmode=disable up