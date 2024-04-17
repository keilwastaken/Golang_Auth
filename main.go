package main

import (
	"Clarity_go/Config"
	"Clarity_go/Helpers"
	"Clarity_go/Server"
	"log"
)

func main() {

	log.Println("Starting Runners App")

	log.Println("Initializing configuration")
	config := Config.Init()

	log.Println("Initializing database")
	mongodb := Server.NewMongoService(config)

	log.Println("Initializing helpers")
	token := Helpers.NewTokenHelper(config)

	log.Println("Initializing HTTP sever")
	httpServer := Server.InitHttpServer(config, mongodb, token)
	httpServer.Start()
}
