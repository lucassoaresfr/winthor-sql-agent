package domain

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Promotion struct {
	CodPrecoProm *int     `json:"codprecoprom,omitempty" gorm:"column:CODPRECOPROM"`
	CodDescricao *int     `json:"coddescricao,omitempty" gorm:"column:CODDESCRICAO"`
	DescPromo    *string  `json:"descpromo,omitempty" gorm:"column:DESCPROMO"`
	CodProd      *int     `json:"codprod,omitempty" gorm:"column:CODPROD"`
	DescProd     *string  `json:"descprod,omitempty" gorm:"column:DESCPROD"`
	PrecoFixo    *string  `json:"precofixo,omitempty" gorm:"column:PRECOFIXO"`
	Estoque      *float64 `json:"estoque,omitempty" gorm:"column:ESTOQUE"`
	Embalagem    *string  `json:"embalagem,omitempty" gorm:"column:EMBALAGEM"`
	Unidade      *string  `json:"unidade,omitempty" gorm:"column:UNIDADE"`
}

type FiltroPromotion struct {
	CodPrecoProm  *int     `form:"codprecoprom"`
	CodDescricao  *int     `form:"coddescricao"`
	DescPromo     string   `form:"descpromo"`
	CodProd       *int     `form:"codprod"`
	DescProd      string   `form:"descprod"`
	Embalagem     string   `form:"embalagem"`
	Unidade       string   `form:"unidade"`
	PrecoFixo     *float64 `form:"precofixo"`
	ApenasEstoque bool     `form:"apenas_estoque"`
	ApenasTipos   bool     `form:"apenas_tipos"` // Quando true, retorna apenas DISTINCT das descrições das ofertas
	OrdenarPor    string   `form:"ordenar_por"`
	Ordem         string   `form:"ordem"`
	Limite        *int     `form:"limite"`
}

// Query base com cálculo correto do estoque disponível no ERP WinThor (ESTOQUE GERENCIAL - RESERVADO - BLOQUEADO)
const QueryPromocaoBase = `SELECT PR.CODPRECOPROM,
       PR.CODDESCRICAO,
       PR.DESCRICAO AS DESCPROMO,
       PR.CODPROD,
       P.DESCRICAO AS DESCPROD,
       TO_CHAR(PR.PRECOFIXO, 'FM999,999,999.00') AS PRECOFIXO,
       (NVL(E.QTESTGER, 0) - NVL(E.QTRESERV, 0) - NVL(E.QTBLOQUEADA, 0)) AS ESTOQUE,
       P.EMBALAGEM,
       P.UNIDADE
FROM PCPRECOPROM PR
INNER JOIN PCPRODUT P ON P.CODPROD = PR.CODPROD
LEFT JOIN PCEST E ON PR.CODPROD = E.CODPROD 
WHERE PR.DTFIMVIGENCIA >= TRUNC(SYSDATE) 
AND E.CODFILIAL = 1`

// Query para buscar apenas os tipos/descrições distintas de promoções ativas
const QueryPromocaoDistinctBase = `SELECT DISTINCT 
       PR.CODDESCRICAO,
       PR.DESCRICAO AS DESCPROMO
FROM PCPRECOPROM PR
INNER JOIN PCPRODUT P ON P.CODPROD = PR.CODPROD
LEFT JOIN PCEST E ON PR.CODPROD = E.CODPROD 
WHERE PR.DTFIMVIGENCIA >= TRUNC(SYSDATE) 
AND E.CODFILIAL = 1`

// ToSQL constrói a instrução SQL com parâmetros nomeados compatíveis com Oracle
func (f FiltroPromotion) ToSQL() (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if f.CodPrecoProm != nil {
		conditions = append(conditions, "PR.CODPRECOPROM = :codprecoprom")
		args = append(args, sql.Named("codprecoprom", *f.CodPrecoProm))
	}
	if f.CodDescricao != nil {
		conditions = append(conditions, "PR.CODDESCRICAO = :coddescricao")
		args = append(args, sql.Named("coddescricao", *f.CodDescricao))
	}
	if strings.TrimSpace(f.DescPromo) != "" {
		conditions = append(conditions, "UPPER(PR.DESCRICAO) LIKE UPPER(:descpromo)")
		args = append(args, sql.Named("descpromo", "%"+strings.TrimSpace(f.DescPromo)+"%"))
	}
	if f.CodProd != nil {
		conditions = append(conditions, "PR.CODPROD = :codprod")
		args = append(args, sql.Named("codprod", *f.CodProd))
	}
	if strings.TrimSpace(f.DescProd) != "" {
		conditions = append(conditions, "UPPER(P.DESCRICAO) LIKE UPPER(:descprod)")
		args = append(args, sql.Named("descprod", "%"+strings.TrimSpace(f.DescProd)+"%"))
	}
	if strings.TrimSpace(f.Embalagem) != "" {
		conditions = append(conditions, "UPPER(P.EMBALAGEM) LIKE UPPER(:embalagem)")
		args = append(args, sql.Named("embalagem", "%"+strings.TrimSpace(f.Embalagem)+"%"))
	}
	if strings.TrimSpace(f.Unidade) != "" {
		conditions = append(conditions, "UPPER(P.UNIDADE) LIKE UPPER(:unidade)")
		args = append(args, sql.Named("unidade", "%"+strings.TrimSpace(f.Unidade)+"%"))
	}
	if f.PrecoFixo != nil {
		conditions = append(conditions, "PR.PRECOFIXO = :precofixo")
		args = append(args, sql.Named("precofixo", *f.PrecoFixo))
	}
	if f.ApenasEstoque {
		conditions = append(conditions, "(NVL(E.QTESTGER, 0) - NVL(E.QTRESERV, 0) - NVL(E.QTBLOQUEADA, 0)) > 0")
	}

	// Seleciona a query base (Normal ou Distinct)
	query := QueryPromocaoBase
	if f.ApenasTipos {
		query = QueryPromocaoDistinctBase
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// Define ordenação padrão conforme a modalidade
	colunaOrdenacao := "PR.CODPRECOPROM"
	if f.ApenasTipos {
		colunaOrdenacao = "PR.DESCRICAO"
	}

	switch strings.ToLower(strings.TrimSpace(f.OrdenarPor)) {
	case "codprod":
		colunaOrdenacao = "PR.CODPROD"
	case "descprod":
		colunaOrdenacao = "P.DESCRICAO"
	case "descpromo":
		colunaOrdenacao = "PR.DESCRICAO"
	case "precofixo":
		colunaOrdenacao = "PR.PRECOFIXO"
	case "estoque":
		colunaOrdenacao = "ESTOQUE"
	case "embalagem":
		colunaOrdenacao = "P.EMBALAGEM"
	case "unidade":
		colunaOrdenacao = "P.UNIDADE"
	case "codprecoprom":
		colunaOrdenacao = "PR.CODPRECOPROM"
	}

	direcao := "DESC"
	if strings.ToUpper(strings.TrimSpace(f.Ordem)) == "ASC" || f.ApenasTipos {
		direcao = "ASC"
	}

	limite := 100
	if f.Limite != nil && *f.Limite > 0 {
		limite = *f.Limite
	}

	query += " ORDER BY " + colunaOrdenacao + " " + direcao

	query = fmt.Sprintf("SELECT * FROM (%s) WHERE ROWNUM <= %d", query, limite)

	return query, args
}

type PromotionRepository interface {
	BuscarPromocoes(ctx context.Context, filtro FiltroPromotion) ([]Promotion, error)
}

type PromotionService interface {
	ListarPromocoes(ctx context.Context, filtro FiltroPromotion) ([]Promotion, error)
}
