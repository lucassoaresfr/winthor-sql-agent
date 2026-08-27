package repository

import (
	"context"

	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
	"gorm.io/gorm"
)

type clientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) domain.ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) BuscarClientes(c context.Context, filtro domain.FiltroCliente) ([]domain.Cliente, error) {
	sqlQuery, args := filtro.ToSQL()

	var clientes []domain.Cliente

	err := r.db.WithContext(c).Raw(sqlQuery, args...).Scan(&clientes).Error
	if err != nil {
		return nil, err
	}

	return clientes, nil
}
