package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() *sql.DB {
	var err error
	// Replace with your MySQL credentials
	dsn := "root:@tcp(127.0.0.1:3306)/preparation?parseTime=true"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	// Test connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
	}

	log.Println("Database connected and initialized!")

	return DB
}
