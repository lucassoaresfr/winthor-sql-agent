package service

import (
	"context"

	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
)

type ProdutosService struct {
	repo domain.ProdutoRepository
}

func NewProdService(repo domain.ProdutoRepository) domain.ProdutoService {
	return &ProdutosService{repo: repo}
}

func (s *ProdutosService) ListarProdutos(c context.Context, filtro domain.FiltroProduto) ([]domain.Produto, error) {
	produtos, err := s.repo.BuscarProdutos(c, filtro)

	if err != nil {
		return nil, err
	}

	return produtos, err
}
