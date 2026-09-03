package domain

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// ItemPedido representa os produtos que compõem o pedido (PCPEDI)
type ItemOrder struct {
	NumPed    int      `json:"numped" gorm:"column:NUMPED"`
	CodProd   int      `json:"codprod" gorm:"column:CODPROD"`
	Descricao *string  `json:"descricao" gorm:"column:DESCRICAO"`
	Qt        float64  `json:"qt" gorm:"column:QT"`
	PVenda    float64  `json:"pvenda" gorm:"column:PVENDA"`
	PTabela   *float64 `json:"ptabela" gorm:"column:PTABELA"`
	PerDesc   *float64 `json:"perdesc" gorm:"column:PERDESC"`
}

// Pedido representa a capa/cabeçalho do pedido de venda (PCPEDC)
type Orders struct {
	NumPed       int         `json:"numped" gorm:"column:NUMPED"`
	Data         *string     `json:"data" gorm:"column:DATA"`
	CodCli       *int        `json:"codcli" gorm:"column:CODCLI"`
	Cliente      *string     `json:"cliente" gorm:"column:CLIENTE"`
	TipoFj       *string     `json:"tipofj" gorm:"TIPOFJ"`
	CodUsur      *int        `json:"codusur" gorm:"column:CODUSUR"`
	NomeVend     *string     `json:"nomevend" gorm:"column:NOMEVEND"`
	Posicao      *string     `json:"posicao" gorm:"column:POSICAO"`
	VlTotal      *float64    `json:"vltotal" gorm:"column:VLTOTAL"`
	VlAtend      *float64    `json:"vlatend" gorm:"column:VLATEND"`
	CodPlPag     *int        `json:"codplpag" gorm:"column:CODPLPAG"`
	DescricaoPag *string     `json:"descricaopag" gorm:"column:DESCRICAOPAG"`
	CodCob       *string     `json:"codcob" gorm:"column:CODCOB"`
	CondVenda    *int        `json:"condvenda" gorm:"column:CONDVENDA"`
	OrigemPed    *string     `json:"origemped" gorm:"column:ORIGEMPED"`
	NumCar       *int        `json:"numcar" gorm:"column:NUMCAR"`
	DtFat        *string     `json:"dtfat" gorm:"column:DTFAT"`
	Obs          *string     `json:"obs" gorm:"column:OBS"`
	Itens        []ItemOrder `json:"itens,omitempty" gorm:"-"`
}

// FiltroPedido mapeia as query string params para a consulta de pedidos
type FiltroOrder struct {
	NumPed     *int   `form:"numped"`
	CodCli     *int   `form:"codcli"`
	Cliente    string `form:"cliente"`
	TipoFj     string `form:"tipofj"`
	CodUsur    *int   `form:"codusur"`
	NomeVend   string `form:"nomevend"`
	Posicao    string `form:"posicao"`
	CodPlPag   *int   `form:"codplpag"`
	CodCob     string `form:"codcob"`
	DtInicio   string `form:"dt_inicio"`
	DtFim      string `form:"dt_fim"`
	OrdenarPor string `form:"ordenar_por"`
	Ordem      string `form:"ordem"`
	Limite     *int   `form:"limite"`
}

const QueryPedidoBase = `SELECT P.NUMPED, TO_CHAR(P.DATA, 'YYYY-MM-DD') AS DATA, P.CODCLI, C.CLIENTE, C.TIPOFJ, 
P.CODUSUR, R.NOME AS NOMEVEND, P.POSICAO, P.VLTOTAL, P.VLATEND, P.CODPLPAG, 
PLAN.DESCRICAO AS DESCRICAOPAG, P.CODCOB, P.CONDVENDA, P.ORIGEMPED, P.NUMCAR, 
TO_CHAR(P.DTFAT, 'YYYY-MM-DD') AS DTFAT, P.OBS
FROM PCPEDC P 
LEFT JOIN PCCLIENT C ON P.CODCLI = C.CODCLI
LEFT JOIN PCUSUARI R ON P.CODUSUR = R.CODUSUR
LEFT JOIN PCPLPAG PLAN ON P.CODPLPAG = PLAN.CODPLPAG 
WHERE 1=1`

const QueryItensPedidoBase = `SELECT I.NUMPED, I.CODPROD, P.DESCRICAO, I.QT, I.PVENDA, I.PTABELA, I.PERDESC
FROM PCPEDI I
LEFT JOIN PCPRODUT P ON I.CODPROD = P.CODPROD
WHERE I.NUMPED = :numped`

// MascararDadosSensiveis mascara o nome do cliente se for Pessoa Física (PF)
func (o *Orders) MascararDadosSensiveis() {
	if o.TipoFj != nil && strings.ToUpper(strings.TrimSpace(*o.TipoFj)) == "F" {
		if o.Cliente != nil && len(*o.Cliente) > 0 {
			nome := strings.TrimSpace(*o.Cliente)
			partes := strings.Fields(nome)
			if len(partes) > 1 {
				// Mantém o primeiro nome e mascara os sobrenomes
				mascarado := partes[0] + " " + strings.Repeat("*", len(nome)-len(partes[0])-1)
				o.Cliente = &mascarado
			} else if len(nome) > 2 {
				// Se for apenas um nome, exibe as 2 primeiras letras e mascara o restante
				mascarado := nome[:2] + strings.Repeat("*", len(nome)-2)
				o.Cliente = &mascarado
			} else {
				mascarado := "***"
				o.Cliente = &mascarado
			}
		}
	}
}

// ToSQL constrói a consulta SQL parametrizada para busca da capa dos pedidos
func (f FiltroOrder) ToSQL() (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if f.NumPed != nil {
		conditions = append(conditions, "P.NUMPED = :numped")
		args = append(args, sql.Named("numped", *f.NumPed))
	}
	if f.CodCli != nil {
		conditions = append(conditions, "P.CODCLI = :codcli")
		args = append(args, sql.Named("codcli", *f.CodCli))
	}
	if strings.TrimSpace(f.Cliente) != "" {
		conditions = append(conditions, "UPPER(C.CLIENTE) LIKE UPPER(:cliente)")
		args = append(args, sql.Named("cliente", "%"+strings.TrimSpace(f.Cliente)+"%"))
	}
	if f.CodUsur != nil {
		conditions = append(conditions, "P.CODUSUR = :codusur")
		args = append(args, sql.Named("codusur", *f.CodUsur))
	}
	if strings.TrimSpace(f.NomeVend) != "" {
		conditions = append(conditions, "UPPER(R.NOME) LIKE UPPER(:nomevend)")
		args = append(args, sql.Named("nomevend", "%"+strings.TrimSpace(f.NomeVend)+"%"))
	}
	if strings.TrimSpace(f.Posicao) != "" {
		conditions = append(conditions, "UPPER(P.POSICAO) = UPPER(:posicao)")
		args = append(args, sql.Named("posicao", strings.TrimSpace(f.Posicao)))
	}
	if f.CodPlPag != nil {
		conditions = append(conditions, "P.CODPLPAG = :codplpag")
		args = append(args, sql.Named("codplpag", *f.CodPlPag))
	}
	if strings.TrimSpace(f.CodCob) != "" {
		conditions = append(conditions, "UPPER(P.CODCOB) = UPPER(:codcob)")
		args = append(args, sql.Named("codcob", strings.TrimSpace(f.CodCob)))
	}
	if strings.TrimSpace(f.DtInicio) != "" {
		conditions = append(conditions, "P.DATA >= TO_DATE(:dt_inicio, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dt_inicio", strings.TrimSpace(f.DtInicio)))
	}
	if strings.TrimSpace(f.DtFim) != "" {
		conditions = append(conditions, "P.DATA <= TO_DATE(:dt_fim, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dt_fim", strings.TrimSpace(f.DtFim)))
	}

	query := QueryPedidoBase
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	colunaOrdenacao := "P.DATA"
	switch strings.ToLower(strings.TrimSpace(f.OrdenarPor)) {
	case "numped":
		colunaOrdenacao = "P.NUMPED"
	case "vltotal":
		colunaOrdenacao = "P.VLTOTAL"
	case "cliente":
		colunaOrdenacao = "C.CLIENTE"
	case "data":
		colunaOrdenacao = "P.DATA"
	}

	direcao := "DESC"
	if strings.ToUpper(strings.TrimSpace(f.Ordem)) == "ASC" {
		direcao = "ASC"
	}

	limite := 50
	if f.Limite != nil && *f.Limite > 0 {
		limite = *f.Limite
	}

	query += " ORDER BY " + colunaOrdenacao + " " + direcao

	query = fmt.Sprintf("SELECT * FROM (%s) WHERE ROWNUM <= %d", query, limite)

	return query, args
}

type OrderRepository interface {
	BuscarPedidos(ctx context.Context, filtro FiltroOrder) ([]Orders, error)
	BuscarItensPedido(ctx context.Context, numPed int) ([]ItemOrder, error)
}

type OrderService interface {
	ListarPedidos(ctx context.Context, filtro FiltroOrder, incluirItens bool) ([]Orders, error)
	ObterPedidoPorNum(ctx context.Context, numPed int) (*Orders, error)
}
