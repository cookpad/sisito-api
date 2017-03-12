package main

import (
	"log"
	"sisito"
)

func main() {
	config, err := sisito.LoadConfig()

	if err != nil {
		log.Fatalf("Load config.tml failed: %s", err)
	}

	driver, err := sisito.NewDriver(config)

	if err != nil {
		log.Fatalf("Create database driver failed: %s", err)
	}

	defer driver.Close()

	server := sisito.NewServer(driver)
	server.Run()
}
