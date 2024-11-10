package main

import (
	"log"

	"github.com/Samuel-Ricardo/CapybaraMQ/internal/infra"
)

func main() {
	cfg := infra.LoadConfig()
	log.Println("Config: ", cfg)
}
