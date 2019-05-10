package main

import (
	"log"

	"github.com/dantin/microservice-go/account"
)

func main() {
	app := account.NewServer(":30000")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
