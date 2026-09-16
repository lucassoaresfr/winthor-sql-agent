package domain

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Finicial representa o lançamento individual de contas a pagar (PCLANC)
type Finicial struct {
	Recnum         int64      `json:"recnum" gorm:"column:RECNUM"`
	NumNota        *int64     `json:"numnota" gorm:"column:NUMNOTA"`
	Duplic         *string    `json:"duplic" gorm:"column:DUPLIC"`
	CodFornec      *int       `json:"codfornec" gorm:"column:CODFORNEC"`
	NomeFornecedor *string    `json:"nome_fornecedor" gorm:"column:NOME_FORNECEDOR"`
	CodConta       *int       `json:"codconta" gorm:"column:CODCONTA"`
	NomeConta      *string    `json:"nome_conta" gorm:"column:NOME_CONTA"`
	Historico      *string    `json:"historico" gorm:"column:HISTORICO"`
	Valor          float64    `json:"valor" gorm:"column:VALOR"`
	VPago          *float64   `json:"vpago" gorm:"column:VPAGO"`
	DtEmissao      *time.Time `json:"dtemissao" gorm:"column:DTEMISSAO"`
	DtLanc         *time.Time `json:"dtlanc" gorm:"column:DTLANC"`
	DtVenc         *time.Time `json:"dtvenc" gorm:"column:DTVENC"`
	DtPagto        *time.Time `json:"dtpagto" gorm:"column:DTPAGTO"`
	StatusTitulo   string     `json:"status_titulo" gorm:"column:STATUS_TITULO"`
	NomeFunc       *string    `json:"nomefunc" gorm:"column:NOMEFUNC"`
}

// PCLancSummary representa a consulta de cálculos e totais agregados de contas a pagar
type PCLancSummary struct {
	QtdTotal      int64   `json:"qtd_total" gorm:"column:QTD_TOTAL"`
	QtdPago       int64   `json:"qtd_pago" gorm:"column:QTD_PAGO"`
	QtdAberto     int64   `json:"qtd_aberto" gorm:"column:QTD_ABERTO"`
	QtdAtrasado   int64   `json:"qtd_atrasado" gorm:"column:QTD_ATRASADO"`
	ValorTotal    float64 `json:"valor_total" gorm:"column:VALOR_TOTAL"`
	ValorPago     float64 `json:"valor_pago" gorm:"column:VALOR_PAGO"`
	ValorAberto   float64 `json:"valor_aberto" gorm:"column:VALOR_ABERTO"`
	ValorAtrasado float64 `json:"valor_atrasado" gorm:"column:VALOR_ATRASADO"`
}

// FilterFinicial define os parâmetros de busca para listagem e cálculos
type FilterFinicial struct {
	Recnum         *int64 `form:"recnum"`
	NumNota        *int64 `form:"numnota"`
	Duplic         string `form:"duplic"`
	CodFornec      *int   `form:"codfornec"`
	NomeFornecedor string `form:"nome_fornecedor"`
	CodConta       *int   `form:"codconta"`
	NomeConta      string `form:"nome_conta"`
	Historico      string `form:"historico"`
	StatusTitulo   string `form:"status_titulo"` // "PAGO", "ATRASADO", "A PAGAR" / "ABERTO"
	NomeFunc       string `form:"nomefunc"`

	// Intervalos de Datas (Formato: YYYY-MM-DD)
	DtLancInicio  string `form:"dtlanc_inicio"`
	DtLancFim     string `form:"dtlanc_fim"`
	DtVencInicio  string `form:"dtvenc_inicio"`
	DtVencFim     string `form:"dtvenc_fim"`
	DtPagtoInicio string `form:"dtpagto_inicio"`
	DtPagtoFim    string `form:"dtpagto_fim"`

	// Intervalos de Valores
	ValorMin *float64 `form:"valor_min"`
	ValorMax *float64 `form:"valor_max"`

	// Ordenação e Paginação
	OrdenarPor string `form:"ordenar_por"`
	Ordem      string `form:"ordem"`
	Limite     *int   `form:"limite"`
}

// CTE base para extrair a massa de dados aplicando os filtros na PCLANC
const QueryLancamentoCTE = `
WITH DADOS AS (
    SELECT 
        L.RECNUM, 
        L.NUMNOTA, 
        L.DUPLIC, 
        L.CODFORNEC, 
        F.FORNECEDOR AS NOME_FORNECEDOR, 
        L.CODCONTA, 
        C.CONTA AS NOME_CONTA, 
        L.HISTORICO, 
        L.VALOR, 
        L.VPAGO, 
        L.DTEMISSAO, 
        L.DTLANC, 
        L.DTVENC, 
        L.DTPAGTO, 
        CASE 
            WHEN L.DTPAGTO IS NOT NULL THEN 'PAGO'
            WHEN L.DTVENC < TRUNC(SYSDATE) THEN 'ATRASADO'
            ELSE 'A PAGAR'
        END AS STATUS_TITULO,
        L.NOMEFUNC
    FROM PCLANC L
        LEFT JOIN PCFORNEC F ON L.CODFORNEC = F.CODFORNEC
        LEFT JOIN PCCONTA C  ON L.CODCONTA = C.CODCONTA
    WHERE 1=1 %s
)`

// BuildConditions monta as cláusulas WHERE baseadas dinamicamente nos parâmetros informados
func (f FilterFinicial) BuildConditions() (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if f.Recnum != nil {
		conditions = append(conditions, "L.RECNUM = :recnum")
		args = append(args, sql.Named("recnum", *f.Recnum))
	}
	if f.NumNota != nil {
		conditions = append(conditions, "L.NUMNOTA = :numnota")
		args = append(args, sql.Named("numnota", *f.NumNota))
	}
	if strings.TrimSpace(f.Duplic) != "" {
		conditions = append(conditions, "UPPER(L.DUPLIC) LIKE UPPER(:duplic)")
		args = append(args, sql.Named("duplic", "%"+strings.TrimSpace(f.Duplic)+"%"))
	}
	if f.CodFornec != nil {
		conditions = append(conditions, "L.CODFORNEC = :codfornec")
		args = append(args, sql.Named("codfornec", *f.CodFornec))
	}
	if strings.TrimSpace(f.NomeFornecedor) != "" {
		conditions = append(conditions, "UPPER(F.FORNECEDOR) LIKE UPPER(:nome_fornecedor)")
		args = append(args, sql.Named("nome_fornecedor", "%"+strings.TrimSpace(f.NomeFornecedor)+"%"))
	}
	if f.CodConta != nil {
		conditions = append(conditions, "L.CODCONTA = :codconta")
		args = append(args, sql.Named("codconta", *f.CodConta))
	}
	if strings.TrimSpace(f.NomeConta) != "" {
		conditions = append(conditions, "UPPER(C.CONTA) LIKE UPPER(:nome_conta)")
		args = append(args, sql.Named("nome_conta", "%"+strings.TrimSpace(f.NomeConta)+"%"))
	}
	if strings.TrimSpace(f.Historico) != "" {
		conditions = append(conditions, "UPPER(L.HISTORICO) LIKE UPPER(:historico)")
		args = append(args, sql.Named("historico", "%"+strings.TrimSpace(f.Historico)+"%"))
	}
	if strings.TrimSpace(f.NomeFunc) != "" {
		conditions = append(conditions, "UPPER(L.NOMEFUNC) LIKE UPPER(:nomefunc)")
		args = append(args, sql.Named("nomefunc", "%"+strings.TrimSpace(f.NomeFunc)+"%"))
	}

	// Filtro por Status do Título (PAGO, ATRASADO, A PAGAR / ABERTO)
	switch strings.ToUpper(strings.TrimSpace(f.StatusTitulo)) {
	case "PAGO":
		conditions = append(conditions, "L.DTPAGTO IS NOT NULL")
	case "ATRASADO":
		conditions = append(conditions, "L.DTPAGTO IS NULL AND L.DTVENC < TRUNC(SYSDATE)")
	case "A PAGAR", "ABERTO":
		conditions = append(conditions, "L.DTPAGTO IS NULL AND L.DTVENC >= TRUNC(SYSDATE)")
	}

	// Filtros de Data de Lançamento
	if strings.TrimSpace(f.DtLancInicio) != "" {
		conditions = append(conditions, "L.DTLANC >= TO_DATE(:dtlanc_inicio, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dtlanc_inicio", strings.TrimSpace(f.DtLancInicio)))
	}
	if strings.TrimSpace(f.DtLancFim) != "" {
		conditions = append(conditions, "L.DTLANC <= TO_DATE(:dtlanc_fim, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dtlanc_fim", strings.TrimSpace(f.DtLancFim)))
	}

	// Filtros de Data de Vencimento
	if strings.TrimSpace(f.DtVencInicio) != "" {
		conditions = append(conditions, "L.DTVENC >= TO_DATE(:dtvenc_inicio, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dtvenc_inicio", strings.TrimSpace(f.DtVencInicio)))
	}
	if strings.TrimSpace(f.DtVencFim) != "" {
		conditions = append(conditions, "L.DTVENC <= TO_DATE(:dtvenc_fim, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dtvenc_fim", strings.TrimSpace(f.DtVencFim)))
	}

	// Filtros de Data de Pagamento
	if strings.TrimSpace(f.DtPagtoInicio) != "" {
		conditions = append(conditions, "L.DTPAGTO >= TO_DATE(:dtpagto_inicio, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dtpagto_inicio", strings.TrimSpace(f.DtPagtoInicio)))
	}
	if strings.TrimSpace(f.DtPagtoFim) != "" {
		conditions = append(conditions, "L.DTPAGTO <= TO_DATE(:dtpagto_fim, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dtpagto_fim", strings.TrimSpace(f.DtPagtoFim)))
	}

	// Filtros de Faixa de Valor
	if f.ValorMin != nil {
		conditions = append(conditions, "L.VALOR >= :valor_min")
		args = append(args, sql.Named("valor_min", *f.ValorMin))
	}
	if f.ValorMax != nil {
		conditions = append(conditions, "L.VALOR <= :valor_max")
		args = append(args, sql.Named("valor_max", *f.ValorMax))
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " AND " + strings.Join(conditions, " AND ")
	}

	return whereClause, args
}

// ToDataSQL gera a instrução SQL com a lista de lançamentos, ordenação e limite via ROWNUM
func (f FilterFinicial) ToDataSQL() (string, []interface{}) {
	whereClause, args := f.BuildConditions()
	cte := fmt.Sprintf(QueryLancamentoCTE, whereClause)

	colunaOrdenacao := "DTVENC"
	switch strings.ToLower(strings.TrimSpace(f.OrdenarPor)) {
	case "dtvenc", "vencimento":
		colunaOrdenacao = "DTVENC"
	case "dtlanc", "lancamento":
		colunaOrdenacao = "DTLANC"
	case "dtpagto", "pagamento":
		colunaOrdenacao = "DTPAGTO"
	case "valor":
		colunaOrdenacao = "VALOR"
	case "numnota", "nota":
		colunaOrdenacao = "NUMNOTA"
	case "fornecedor":
		colunaOrdenacao = "NOME_FORNECEDOR"
	}

	direcao := "DESC NULLS LAST"
	if strings.ToUpper(strings.TrimSpace(f.Ordem)) == "ASC" {
		direcao = "ASC NULLS LAST"
	}

	limite := 50
	if f.Limite != nil && *f.Limite > 0 {
		limite = *f.Limite
	}

	query := fmt.Sprintf("%s\nSELECT * FROM DADOS ORDER BY %s %s", cte, colunaOrdenacao, direcao)
	fullQuery := fmt.Sprintf("SELECT * FROM (%s) WHERE ROWNUM <= %d", query, limite)

	return fullQuery, args
}

// ToCalculosSQL gera a instrução SQL para os Totais Agregados respeitando o limite do ROWNUM
func (f FilterFinicial) ToCalculosSQL() (string, []interface{}) {
	whereClause, args := f.BuildConditions()
	cte := fmt.Sprintf(QueryLancamentoCTE, whereClause)

	colunaOrdenacao := "DTVENC"
	switch strings.ToLower(strings.TrimSpace(f.OrdenarPor)) {
	case "dtvenc", "vencimento":
		colunaOrdenacao = "DTVENC"
	case "dtlanc", "lancamento":
		colunaOrdenacao = "DTLANC"
	case "dtpagto", "pagamento":
		colunaOrdenacao = "DTPAGTO"
	case "valor":
		colunaOrdenacao = "VALOR"
	case "numnota", "nota":
		colunaOrdenacao = "NUMNOTA"
	case "fornecedor":
		colunaOrdenacao = "NOME_FORNECEDOR"
	}

	direcao := "DESC NULLS LAST"
	if strings.ToUpper(strings.TrimSpace(f.Ordem)) == "ASC" {
		direcao = "ASC NULLS LAST"
	}

	limite := 50
	if f.Limite != nil && *f.Limite > 0 {
		limite = *f.Limite
	}

	summaryQuery := fmt.Sprintf(`%s
SELECT
    COUNT(*) AS QTD_TOTAL,
    COUNT(CASE WHEN DTPAGTO IS NOT NULL THEN 1 END) AS QTD_PAGO,
    COUNT(CASE WHEN DTPAGTO IS NULL AND DTVENC >= TRUNC(SYSDATE) THEN 1 END) AS QTD_ABERTO,
    COUNT(CASE WHEN DTPAGTO IS NULL AND DTVENC < TRUNC(SYSDATE) THEN 1 END) AS QTD_ATRASADO,
    NVL(SUM(VALOR), 0) AS VALOR_TOTAL,
    NVL(SUM(VPAGO), 0) AS VALOR_PAGO,
    NVL(SUM(CASE WHEN DTPAGTO IS NULL AND DTVENC >= TRUNC(SYSDATE) THEN VALOR ELSE 0 END), 0) AS VALOR_ABERTO,
    NVL(SUM(CASE WHEN DTPAGTO IS NULL AND DTVENC < TRUNC(SYSDATE) THEN VALOR ELSE 0 END), 0) AS VALOR_ATRASADO
FROM (
    SELECT D.* FROM DADOS D ORDER BY %s %s
) WHERE ROWNUM <= %d`, cte, colunaOrdenacao, direcao, limite)

	return summaryQuery, args
}

type FinicialRepository interface {
	SearchLaunch(c context.Context, filtro FilterFinicial) ([]Finicial, error)
	GetPCLancCalculos(c context.Context, filtro FilterFinicial) (*PCLancSummary, error)
}

type FinicialService interface {
	ListLaunch(c context.Context, filtro FilterFinicial) ([]Finicial, *PCLancSummary, error)
}
