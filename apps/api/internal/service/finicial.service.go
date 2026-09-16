package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
)

type finicialService struct {
	repo domain.FinicialRepository
}

func NewFinicialService(repo domain.FinicialRepository) domain.FinicialService {
	return &finicialService{repo: repo}
}

func (s finicialService) ListLaunch(c context.Context, filter domain.FilterFinicial) ([]domain.Finicial, *domain.PCLancSummary, error) {
	var (
		launch              []domain.Finicial
		summary             *domain.PCLancSummary
		errData, errSummary error
		wg                  sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		launch, errData = s.repo.SearchLaunch(c, filter)
	}()

	go func() {
		defer wg.Done()
		summary, errSummary = s.repo.GetPCLancCalculos(c, filter)
	}()

	wg.Wait()

	if errData != nil {
		return nil, nil, fmt.Errorf("erro ao buscar lançamentos: %w", errData)
	}

	if errSummary != nil {
		return nil, nil, fmt.Errorf("erro ao buscar cálculos de resumo: %w", errSummary)
	}

	return launch, summary, nil
}
