package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3308)/sys?")
	fmt.Print(err)
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	fmt.Print(err)
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/mysql",
		"mysql",
		driver,
	)
	fmt.Print(err)

	err = m.Steps(2)
	fmt.Print(err)
}
