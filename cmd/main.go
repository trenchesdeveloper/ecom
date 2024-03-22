package main

import (
	"github.com/trenchesdeveloper/go-ecom/cmd/api"
	"log"
)

func main() {
	server := api.NewApplication(":4000", nil)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
