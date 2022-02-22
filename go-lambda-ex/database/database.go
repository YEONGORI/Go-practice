package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("host")
	password = os.Getenv("password")
	port     = 5432
	user     = "postgres"
	database = "postgres"
)

func GetConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, database)

	// Open DB를 연다. 첫번째 인자로는 해당 DB driver(pq: postgres driver)의 이름을 받고
	// 두번째 인자로는 Connection Information을 받는다.
	return sql.Open("postgres", psqlInfo)
}

const createEmployeesTableSQL = `
CREATE TABLE employees(
	id serial,
	name varchar,
	first_name varchar,
	last_name varchar
);
`

func CreateEmployeesTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, createEmployeesTableSQL)
	return err
}
