package repository

import (
	"context"

	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
	"gorm.io/gorm"
)

type produtosRepository struct {
	db *gorm.DB
}

func NewProdRepository(db *gorm.DB) domain.ProdutoRepository {
	return &produtosRepository{db: db}
}

func (r *produtosRepository) BuscarProdutos(c context.Context, filtro domain.FiltroProduto) ([]domain.Produto, error) {
	sqlQuery, args := filtro.ToSQL()

	var produtos []domain.Produto

	err := r.db.WithContext(c).Raw(sqlQuery, args...).Scan(&produtos).Error

	if err != nil {
		return nil, err
	}

	return produtos, nil
}
