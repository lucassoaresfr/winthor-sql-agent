package repository

import (
	"context"

	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
	"gorm.io/gorm"
)

type promotionRepository struct {
	db *gorm.DB
}

func NewPromotionRepository(db *gorm.DB) domain.PromotionRepository {
	return &promotionRepository{db: db}
}

func (r *promotionRepository) BuscarPromocoes(c context.Context, filtro domain.FiltroPromotion) ([]domain.Promotion, error) {
	sqlQuery, args := filtro.ToSQL()

	var promotions []domain.Promotion

	err := r.db.WithContext(c).Raw(sqlQuery, args...).Scan(&promotions).Error
	if err != nil {
		return nil, err
	}

	return promotions, nil
}
