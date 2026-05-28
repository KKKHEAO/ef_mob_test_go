package main

import (
	"ef_mob_test_go/config"
	"log"
	"os"
)

func main() {
	log.Println("Service started")
	_, err := config.GetConfigByFilename(os.Getenv("config"))
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

}
