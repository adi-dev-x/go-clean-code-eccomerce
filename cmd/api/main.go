package main

import (
	"fmt"

	"log"
	"myproject/pkg/config"
	db "myproject/pkg/database"
	"myproject/pkg/di"
)

func main() {
	// Load configuration
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	errs := db.InitRedis()
	if errs != nil {
		fmt.Printf("Error initializing Redis: %s\n", err.Error())
	} else {
		fmt.Println("Redis connection successful!")
	}

	server, err := di.InitializeEvent(conf)
	if err != nil {
		log.Fatal("failed to initialize the files")
	}

	server.Start(conf)
}
