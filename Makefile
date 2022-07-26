postgres:
	docker run --name=db -ti -e POSTGRES_PASSWORD='root' -ti -p 5436:5432 -d --rm postgres:12.11

postgres_migrate:
	migrate -source file://migrations/postgres -database postgres://postgres:root@localhost:5436/postgres?sslmode=disable up

mysql:
	 docker run --name db -e MYSQL_ROOT_PASSWORD='root' -p 3308:3306 -d --rm  mysql:8.0

mysql_migrate:
	migrate -source file://migrations/mysql -database mysql://root:root@tcp(localhost:3308)/sys up
