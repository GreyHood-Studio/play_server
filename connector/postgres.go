package connector

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type PostgresConn struct {
	client	*sql.DB
	addr		string
	user		string
	pwd		string
	db		string
}

// Postgresql Connect 로직
func (conn *PostgresConn) connectPostgres() {
	//"postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=verify-full", conn.user, conn.pwd, conn.db, conn.db)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		// error 로그
	}
	conn.client = db
}