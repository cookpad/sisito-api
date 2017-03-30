package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"sisito"
)

func main() {
	flags := sisito.ParseFlag()

	config, err := sisito.LoadConfig(flags)

	if err != nil {
		log.Fatalf("Load config.tml failed: %s", err)
	}

	var out io.Writer

	if config.Server.Log != "" {
		file, err := os.OpenFile(config.Server.Log, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()
		out = file
	} else {
		out = os.Stdout
	}

	driver, err := sisito.NewDriver(config, gin.Mode() == "debug", out)

	if err != nil {
		log.Fatalf("Create database driver failed: %s", err)
	}

	defer driver.Close()

	server := sisito.NewServer(config, driver, out)
	server.Run()
}
