package main

import (
	"fmt"
	"log"

	"cuboid-challenge/app/config"
	"cuboid-challenge/app/db"
	"cuboid-challenge/app/router"
)

func main() {
	config.Load()
	db.Connect()

	r := router.Setup()
	addr := fmt.Sprintf(":%s", config.ENV.Port)

	if err := r.Run(addr); err != nil {
		log.Fatalln("Failed to start the application.", err)
	}
}
