package main

import "database/sql"

func OpenDB() *sql.DB {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/dbname?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func CloseDB() {
	if err := DB.Close(); err != nil {
		panic(err.Error())
	}
}
