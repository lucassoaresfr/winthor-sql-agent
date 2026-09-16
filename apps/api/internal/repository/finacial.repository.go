package repository

import (
	"context"

	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
	"gorm.io/gorm"
)

type finicialRepository struct {
	db *gorm.DB
}

func NewFinicialRepository(db *gorm.DB) domain.FinicialRepository {
	return &finicialRepository{db: db}
}

func (r *finicialRepository) SearchLaunch(c context.Context, filter domain.FilterFinicial) ([]domain.Finicial, error) {
	sqlQuery, args := filter.ToDataSQL()

	var launch []domain.Finicial

	err := r.db.WithContext(c).Raw(sqlQuery, args...).Scan(&launch).Error

	if err != nil {
		return nil, err
	}

	return launch, nil
}

func (r *finicialRepository) GetPCLancCalculos(c context.Context, filter domain.FilterFinicial) (*domain.PCLancSummary, error) {
	sqlQuery, args := filter.ToCalculosSQL()

	var sumary domain.PCLancSummary

	err := r.db.WithContext(c).Raw(sqlQuery, args...).Scan(&sumary).Error

	if err != nil {
		return nil, err
	}

	return &sumary, nil
}
