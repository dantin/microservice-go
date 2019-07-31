package main

import (
	"log"

	"github.com/dantin/microservice-go/account/service"
)

func main() {
	app := service.NewServer(":6767")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
