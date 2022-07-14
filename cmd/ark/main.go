package main

import (
	"fmt"
	"log"

	"github.com/alwindoss/ark"
	"github.com/alwindoss/ark/internal/engine"
	"github.com/caarlos0/env/v6"
)

func main() {
	fmt.Println("Welcome to the Ark")
	cfg := ark.Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", cfg)
	log.Fatal(engine.Run(&cfg))
}
