package main

import (
	"log"
	"sisito"
)

func main() {
	flags := sisito.ParseFlag()

	config, err := sisito.LoadConfig(flags)

	if err != nil {
		log.Fatalf("Load config.tml failed: %s", err)
	}

	driver, err := sisito.NewDriver(config)

	if err != nil {
		log.Fatalf("Create database driver failed: %s", err)
	}

	defer driver.Close()

	server := sisito.NewServer(config, driver)
	server.Run()
}
