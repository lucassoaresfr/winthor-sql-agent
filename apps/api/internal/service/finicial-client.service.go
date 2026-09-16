package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
)

type pcPrestService struct {
	repo domain.PCPrestRepository
}

func NewPCPresService(repo domain.PCPrestRepository) domain.PCPrestService {
	return &pcPrestService{repo: repo}
}

func (s *pcPrestService) ListPCPrest(c context.Context, filter domain.FilterPCPrest) ([]domain.PCPrest, *domain.PCPrestSummary, error) {
	var (
		titulos             []domain.PCPrest
		resumo              *domain.PCPrestSummary
		errData, errSummary error
		wg                  sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		titulos, errData = s.repo.SearchPCPrest(c, filter)
	}()

	go func() {
		defer wg.Done()
		resumo, errSummary = s.repo.GetPCPrestCalculos(c, filter)
	}()

	wg.Wait()

	if errData != nil {
		return nil, nil, fmt.Errorf("falha ao carregar listagem: %w", errData)
	}
	if errData != nil {
		return nil, nil, fmt.Errorf("falha ao carregar resumo de valores: %w", errSummary)
	}

	return titulos, resumo, nil
}
