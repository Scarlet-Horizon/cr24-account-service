package main

import (
	"log"
	"main/env"
)

func main() {
	err := env.Load("env/.env")
	if err != nil {
		log.Fatal(err)
	}
}
