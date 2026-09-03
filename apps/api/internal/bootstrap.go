package internal

import (
	"fmt"
	"log"

	"github.com/lucassoaresfr/winthor-sql-agent.git/config"
)

func Start() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("erro ao carregar configurações: %w", err)
	}

	db, err := config.NewOracleDatabase(cfg.DB)
	if err != nil {
		return fmt.Errorf("erro ao conectar no banco Oracle: %w", err)
	}

	dbPostgres, err := config.NewPostgresDB(&cfg.Postgres)
	if err != nil {
		return fmt.Errorf("erro ao conectar no banco PostgreSQL: %w", err)
	}

	// Passando ambos os bancos de dados
	handlers := RegisterDependencies(db, dbPostgres)

	router := NewRouter(handlers, cfg)

	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 Servidor HTTP rodando na porta %s...", cfg.Port)

	return router.Run(serverAddr)
}
