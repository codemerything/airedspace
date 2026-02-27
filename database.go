package main

import (
	"database/sql"
	"log"
)

type DB struct {
	*sql.DB
}

func NewConnection(connStr string) (*DB, error) {
	conn, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")

	return &DB{
		DB: conn,
	}, nil
}
