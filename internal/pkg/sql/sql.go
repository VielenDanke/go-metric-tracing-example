package sql

import (
	"database/sql"
)

func OpenSQLConnection(url string) (*sql.DB, error) {
	conn, connErr := sql.Open("postgres", url)

	if connErr != nil {
		return nil, connErr
	}
	if pingErr := conn.Ping(); pingErr != nil {
		return nil, pingErr
	}
	return conn, nil
}
