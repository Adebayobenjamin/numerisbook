package main

import (
	"log"

	"github.com/Adebayobenjamin/numerisbook/pkg/configs"
)

func init() {
	configs.NewEnvironment()
}

func main() {
	server, err := BootstrapServer()
	if err != nil {
		log.Fatalf("Failed to bootstrap server: %v", err)
	}
	server.Start()
}
