package main

import (
	"errors"
	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/trenchesdeveloper/go-ecom/config"
	"github.com/trenchesdeveloper/go-ecom/db"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := db.NewDB(mysqlCfg.Config{
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

	driver, err := mysql.WithInstance(db, &mysql.Config{})

	if err != nil {
		log.Fatal(err)

	}

	m, err := migrate.NewWithDatabaseInstance("file://cmd/migrate/migrations", "mysql", driver)

	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[(len(os.Args) - 1)]

	if cmd == "up" {
		err = m.Up()

		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
	}

	if cmd == "down" {
		err = m.Down()

		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Migration successful")
}
