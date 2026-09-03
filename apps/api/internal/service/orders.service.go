package service

import (
	"context"

	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
)

type orderService struct {
	repo domain.OrderRepository
}

func NewOrderService(repo domain.OrderRepository) domain.OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) ListarPedidos(c context.Context, filtro domain.FiltroOrder, incluirItems bool) ([]domain.Orders, error) {
	orders, err := s.repo.BuscarPedidos(c, filtro)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		orders[i].MascararDadosSensiveis()

		if incluirItems {
			items, err := s.repo.BuscarItensPedido(c, orders[i].NumPed)
			if err == nil {
				orders[i].Itens = items
			}
		}
	}

	return orders, nil
}

func (s *orderService) ObterPedidoPorNum(ctx context.Context, numPed int) (*domain.Orders, error) {
	filtro := domain.FiltroOrder{
		NumPed: &numPed,
	}

	orders, err := s.repo.BuscarPedidos(ctx, filtro)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, nil
	}

	order := orders[0]

	items, err := s.repo.BuscarItensPedido(ctx, numPed)
	if err != nil {
		return nil, err
	}

	order.Itens = items

	return &order, nil
}
