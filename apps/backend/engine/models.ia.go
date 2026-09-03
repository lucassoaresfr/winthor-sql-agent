package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var ModelosFallback = []string{
	"gemini-3.5-flash-lite",
	"gemini-3.1-flash-lite",
	"gemini-3-flash-preview",
	"gemini-3.5-flash",
	"gemini-2.5-flash",
	"gemini-2.5-flash-lite",
}

// var ModelosFallback = []string{
// 	"gemini-2.5-flash",      // Prioridade 1: Mais rápido e econômico
// 	"gemini-2.0-flash",      // Prioridade 2: Backup rápido
// 	"gemini-2.5-pro",        // Prioridade 3: Raciocínio profundo
// 	"gemini-2.0-flash-lite", // Prioridade 4: Consumo baixo de cota
// 	"gemini-1.5-flash",      // Prioridade 5: Backup legado estável
// }

type ModelManager struct {
	rdb *redis.Client
}

func NewModelManager(rdb *redis.Client) *ModelManager {
	return &ModelManager{rdb: rdb}
}

// ObterModeloDisponivel verifica se a chave existe usando WithContext
func (m *ModelManager) ObterModeloDisponivel(c context.Context) string {
	for _, modelo := range ModelosFallback {
		key := fmt.Sprintf("model:cooldown:%s", modelo)

		// Adicionado .WithContext(c) e passado apenas a key no Exists
		exists, err := m.rdb.WithContext(c).Exists(key).Result()
		if err == nil && exists > 0 {
			continue
		}

		return modelo
	}

	return ModelosFallback[0]
}

func (m *ModelManager) RegistrarCooldownBloqueio(c context.Context, mod string, minutos int) error {
	key := fmt.Sprintf("model:cooldown:%s", mod)
	duracao := time.Duration(minutos) * time.Minute

	return m.rdb.WithContext(c).Set(key, "rate_limit_exceeded", duracao).Err()
}

// EstaEmCooldown verifica se um modelo específico possui chave ativa no Redis
func (m *ModelManager) EstaEmCooldown(c context.Context, modelo string) bool {
	key := fmt.Sprintf("model:cooldown:%s", modelo)
	exists, err := m.rdb.WithContext(c).Exists(key).Result()
	return err == nil && exists > 0
}
