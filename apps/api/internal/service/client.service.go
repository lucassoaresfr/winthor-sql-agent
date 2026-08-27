package service

import (
	"context"

	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
)

type ClientService struct {
	repo domain.ClientRepository
}

func NewClientService(repo domain.ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) ListarClientes(ctx context.Context, filtro domain.FiltroCliente) ([]domain.Cliente, error) {
	clientes, err := s.repo.BuscarClientes(ctx, filtro)
	if err != nil {
		return nil, err
	}

	// Aplica a política LGPD de mascaramento chamando o método da struct de domínio
	for i := range clientes {
		clientes[i].MascararDadosSensiveis()
	}

	return clientes, nil
}
