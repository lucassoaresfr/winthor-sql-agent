package service

import (
	"context"

	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
)

type promotionService struct {
	repo domain.PromotionRepository
}

func NewPromotionService(repo domain.PromotionRepository) domain.PromotionService {
	return &promotionService{repo: repo}
}

func (s *promotionService) ListarPromocoes(c context.Context, filtro domain.FiltroPromotion) ([]domain.Promotion, error) {
	promotion, err := s.repo.BuscarPromocoes(c, filtro)
	if err != nil {
		return nil, err
	}

	return promotion, nil
}
