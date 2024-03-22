package main

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/trenchesdeveloper/go-ecom/cmd/api"
	"github.com/trenchesdeveloper/go-ecom/config"
	"github.com/trenchesdeveloper/go-ecom/db"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := db.NewDB(mysql.Config{
		User: config.Envs.DBUser,
		//Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)

	}
	server := api.NewApplication(fmt.Sprintf(":%s", config.Envs.Port), db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
