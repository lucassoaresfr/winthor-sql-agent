package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	go_ora "github.com/sijms/go-ora/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	ServiceName string
}

func NewOracleDatabase(cfg DBConfig) (*gorm.DB, error) {
	// 1. Monta a DSN do Oracle usando go-ora
	dsn := go_ora.BuildUrl(
		cfg.Host,
		cfg.Port,
		cfg.ServiceName,
		cfg.User,
		cfg.Password,
		map[string]string{"TIMEOUT": "30"},
	)

	// 2. Conecta via database/sql nativo com driver "oracle" (go-ora)
	sqlDB, err := sql.Open("oracle", dsn)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir conexão sql com oracle: %w", err)
	}

	// 3. Pool de Conexões ajustado para o Oracle 11g
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("falha no ping do oracle: %w", err)
	}

	// 4. Inicializa o GORM ignorando a checagem automática de versão do MySQL
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true, // <--- AQUI: Impede o GORM de disparar SELECT VERSION() sem FROM DUAL
	}), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return nil, fmt.Errorf("falha ao inicializar gorm com oracle: %w", err)
	}

	log.Println("Conexão com Oracle 11g (go-ora) estabelecida com sucesso!")
	return gormDB, nil
}
