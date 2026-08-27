package main

import (
	"log"

	"github.com/lucassoaresfr/winthor-sql-agent.git/internal"
)

func main() {
	if err := internal.Start(); err != nil {
		log.Fatalf("Falha crítica ao iniciar a aplicação: %v", err)
	}
}
