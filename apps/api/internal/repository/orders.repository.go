package repository

import (
	"context"
	"database/sql"

	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrdersRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) BuscarPedidos(c context.Context, filtro domain.FiltroOrder) ([]domain.Orders, error) {
	query, args := filtro.ToSQL()
	var orders []domain.Orders

	err := r.db.WithContext(c).Raw(query, args...).Scan(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) BuscarItensPedido(ctx context.Context, numped int) ([]domain.ItemOrder, error) {
	var items []domain.ItemOrder
	arg := sql.Named("numped", numped)

	err := r.db.WithContext(ctx).Raw(domain.QueryItensPedidoBase, arg).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
