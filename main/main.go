package main

import (
	"database/sql"
	"log"
	"slice/main/repositories"
	"slice/main/routes"
	"slice/main/services"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/preparation")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	repo := repositories.NewWalletRepository(db)
	service := services.NewWalletService(repo)
	router := routes.SetupRouter(service)

	router.Run(":8080")
}
