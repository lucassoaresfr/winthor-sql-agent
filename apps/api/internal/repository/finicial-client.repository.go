package repository

import (
	"context"
	"errors"

	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
	"gorm.io/gorm"
)

type pcPrestRepository struct {
	db *gorm.DB
}

func NewFinicialClientRepository(db *gorm.DB) domain.PCPrestRepository {
	return &pcPrestRepository{db: db}
}

func (r *pcPrestRepository) SearchPCPrest(c context.Context, filter domain.FilterPCPrest) ([]domain.PCPrest, error) {
	sqlQuery, args := filter.ToDataSQL()

	var titulos []domain.PCPrest

	err := r.db.WithContext(c).Raw(sqlQuery, args...).Scan(&titulos).Error

	if err != nil {
		return nil, err
	}

	return titulos, nil
}

func (r *pcPrestRepository) GetPCPrestCalculos(ctx context.Context, filtro domain.FilterPCPrest) (*domain.PCPrestSummary, error) {
	sqlQuery, args := filtro.ToCalculosSQL()

	var summary domain.PCPrestSummary

	err := r.db.WithContext(ctx).Raw(sqlQuery, args...).Scan(&summary).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &domain.PCPrestSummary{}, nil
		}
		return nil, err
	}

	return &summary, nil
}