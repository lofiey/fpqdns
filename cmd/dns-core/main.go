package main

import (
	"log"

	"fpqdns/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
