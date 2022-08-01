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

func SetupData(conn *sql.DB) error {
	_, err := conn.Exec("drop table if exists users")
	_, err = conn.Exec("create table users(id bigserial primary key, username varchar, email varchar)")
	_, err = conn.Exec("insert into users(username, email) values ('Vlad', 'vlad@mail.ru')")
	_, err = conn.Exec("insert into users(username, email) values ('Vika', 'vika@mail.ru')")
	return err
}
