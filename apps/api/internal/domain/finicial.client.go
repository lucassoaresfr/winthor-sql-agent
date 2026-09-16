package domain

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// PCPrest representa o título individual de contas a receber com todos os campos calculados e originais do SQL
type PCPrest struct {
	CodCli        *int       `json:"codcli" gorm:"column:CODCLI"`
	Cliente       *string    `json:"cliente" gorm:"column:CLIENTE"`
	NumTransVenda *int64     `json:"numtransvenda" gorm:"column:NUMTRANSVENDA"`
	Duplic        *int64     `json:"duplic" gorm:"column:DUPLIC"`
	Prest         *string    `json:"prest" gorm:"column:PREST"`
	DtEmissao     *time.Time `json:"dtemissao" gorm:"column:DTEMISSAO"`
	DtVenc        *time.Time `json:"dtvenc" gorm:"column:DTVENC"`
	DtVencOrig    *time.Time `json:"dtvencorig" gorm:"column:DTVENCORIG"`
	Valor         float64    `json:"valor" gorm:"column:VALOR"`
	CodCob        *string    `json:"codcob" gorm:"column:CODCOB"`
	CodCobOrig    *string    `json:"codcoborig" gorm:"column:CODCOBORIG"`
	CodUsur       *int       `json:"codusur" gorm:"column:CODUSUR"`
	NomeRCA       *string    `json:"nome_rca" gorm:"column:RCA"`
	DtBaixa       *time.Time `json:"dtbaixa" gorm:"column:DTBAIXA"`
	DtPag         *time.Time `json:"dtpag" gorm:"column:DTPAG"`
	VPago         *float64   `json:"vpago" gorm:"column:VPAGO"`
	TxJuros       *float64   `json:"txjuros" gorm:"column:TXJUROS"`
	PercMulta     float64    `json:"percmulta" gorm:"column:PERCMULTA"`
	DiasAtraso    int64      `json:"dias_atraso" gorm:"column:DIASATRASO"`
	VlJurosDia    float64    `json:"vl_juros_dia" gorm:"column:VLJUROSDIA"`
	JurosAtraso   float64    `json:"juros_atraso" gorm:"column:JUROSATRASO"`
	VlMulta       float64    `json:"vl_multa" gorm:"column:VLMULTA"`
	VlTotal       float64    `json:"vl_total" gorm:"column:VLTOTAL"`
}

// PCPrestSummary representa a consulta de cálculos e totais agregados
type PCPrestSummary struct {
	QtdTotal      int64   `json:"qtd_total" gorm:"column:QTD_TOTAL"`
	QtdPago       int64   `json:"qtd_pago" gorm:"column:QTD_PAGO"`
	QtdAberto     int64   `json:"qtd_aberto" gorm:"column:QTD_ABERTO"`
	QtdAtrasado   int64   `json:"qtd_atrasado" gorm:"column:QTD_ATRASADO"`
	ValorTotal    float64 `json:"valor_total" gorm:"column:VALOR_TOTAL"`
	ValorPago     float64 `json:"valor_pago" gorm:"column:VALOR_PAGO"`
	ValorAberto   float64 `json:"valor_aberto" gorm:"column:VALOR_ABERTO"`
	ValorAtrasado float64 `json:"valor_atrasado" gorm:"column:VALOR_ATRASADO"`
}

// FilterPCPrest define os parâmetros de busca para listagem e cálculos
type FilterPCPrest struct {
	NumTransVenda *int64 `form:"numtransvenda"`
	Duplic        *int64 `form:"duplic"`
	Prest         string `form:"prest"`
	CodCli        *int   `form:"codcli"`
	Cliente       string `form:"cliente"`
	CodUsur       *int   `form:"codusur"`
	CodCob        string `form:"codcob"`
	StatusTitulo  string `form:"status_titulo"` // "PAGO", "ATRASADO", "ABERTO" / "A RECEBER"

	// Intervalos de Datas (Formato: YYYY-MM-DD)
	DtEmissaoInicio string `form:"dtemissao_inicio"`
	DtEmissaoFim    string `form:"dtemissao_fim"`
	DtVencInicio    string `form:"dtvenc_inicio"`
	DtVencFim       string `form:"dtvenc_fim"`
	DtPagInicio     string `form:"dtpag_inicio"`
	DtPagFim        string `form:"dtpag_fim"`

	// Intervalos de Valores
	ValorMin *float64 `form:"valor_min"`
	ValorMax *float64 `form:"valor_max"`

	// Ordenação e Paginação
	OrdenarPor string `form:"ordenar_por"`
	Ordem      string `form:"ordem"`
	Limite     *int   `form:"limite"`
}

// CTE com o SELECT base que extrai a massa de dados do Oracle (Usa LEFT JOINs para evitar ocultar registros)
const QueryPCPrestCTE = `
WITH DADOS AS (
    SELECT
        P.CODCLI,
        C.CLIENTE,
        P.NUMTRANSVENDA,
        P.DUPLIC,
        P.PREST,
        P.DTEMISSAO,
        P.DTVENC,
        P.DTVENCORIG,
        P.VALOR,
        P.CODCOB,
        P.CODCOBORIG,
        P.CODUSUR,
        U.NOME AS RCA,
        P.DTBAIXA,
        P.DTPAG,
        P.VPAGO,
        B.TXJUROS,
        NVL(B.PERCMULTA, 0) AS PERCMULTA,
        CASE
            WHEN ROUND(TRUNC(SYSDATE) - P.DTVENC, 0) < 0 OR P.DTPAG < P.DTVENC THEN 0
            WHEN P.DTPAG >= P.DTVENC THEN ROUND(P.DTPAG - P.DTVENC, 0)
            ELSE ROUND(TRUNC(SYSDATE) - P.DTVENC, 0)
        END AS DIASATRASO
    FROM PCPREST P
    JOIN PCCLIENT C ON C.CODCLI = P.CODCLI
    LEFT JOIN PCCOB B ON B.CODCOB = P.CODCOB
    LEFT JOIN PCUSUARI U ON U.CODUSUR = P.CODUSUR
    LEFT JOIN PCSUPERV S ON S.CODSUPERVISOR = U.CODSUPERVISOR
    WHERE 1=1 %s
)`

// Query final para obter todas as colunas mapeadas
const QueryPCPrestCalculada = `
SELECT
    CODCLI,
    CLIENTE,
    NUMTRANSVENDA,
    DUPLIC,
    PREST,
    TRUNC(DTEMISSAO) AS DTEMISSAO,
    DTVENC,
    DTVENCORIG,
    VALOR,
    CODCOB,
    CODCOBORIG,
    DIASATRASO,
    CASE
        WHEN CODCOB IN ('DEVT','DEVP') OR DTVENC > TRUNC(SYSDATE) OR DIASATRASO = 0 THEN 0
        ELSE ROUND(VALOR * ((TXJUROS / 100) / 30), 2)
    END AS VLJUROSDIA,
    CASE
        WHEN DTVENC > TRUNC(SYSDATE) OR DIASATRASO = 0 THEN 0
        ELSE ROUND((VALOR * ((TXJUROS / 100) / 30)) * DIASATRASO, 2)
    END AS JUROSATRASO,
    PERCMULTA,
    CASE
        WHEN CODCOB IN ('DEVT','DEVP') OR DTVENC > TRUNC(SYSDATE) OR DIASATRASO = 0 THEN 0
        ELSE ROUND(VALOR * (PERCMULTA / 100), 2)
    END AS VLMULTA,
    VALOR
    + CASE
        WHEN DTVENC > TRUNC(SYSDATE) OR DIASATRASO = 0 THEN 0
        ELSE ROUND((VALOR * ((TXJUROS / 100) / 30)) * DIASATRASO, 2)
      END
    + CASE
        WHEN CODCOB IN ('DEVT','DEVP') OR DTVENC > TRUNC(SYSDATE) OR DIASATRASO = 0 THEN 0
        ELSE ROUND(VALOR * (PERCMULTA / 100), 2)
      END AS VLTOTAL,
    CODUSUR,
    RCA,
    DTBAIXA,
    DTPAG,
    VPAGO,
    TXJUROS
FROM DADOS`

// BuildConditions monta as cláusulas WHERE baseadas apenas nos parâmetros informados
func (f FilterPCPrest) BuildConditions() (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if f.NumTransVenda != nil {
		conditions = append(conditions, "P.NUMTRANSVENDA = :numtransvenda")
		args = append(args, sql.Named("numtransvenda", *f.NumTransVenda))
	}
	if f.Duplic != nil {
		conditions = append(conditions, "P.DUPLIC = :duplic")
		args = append(args, sql.Named("duplic", *f.Duplic))
	}
	if strings.TrimSpace(f.Prest) != "" {
		conditions = append(conditions, "UPPER(P.PREST) = UPPER(:prest)")
		args = append(args, sql.Named("prest", strings.TrimSpace(f.Prest)))
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
	if strings.TrimSpace(f.CodCob) != "" {
		conditions = append(conditions, "UPPER(P.CODCOB) = UPPER(:codcob)")
		args = append(args, sql.Named("codcob", strings.TrimSpace(f.CodCob)))
	}

	// Filtro de Status
	switch strings.ToUpper(strings.TrimSpace(f.StatusTitulo)) {
	case "PAGO":
		conditions = append(conditions, "P.DTPAG IS NOT NULL")
	case "ATRASADO":
		conditions = append(conditions, "P.DTPAG IS NULL AND P.DTVENC < TRUNC(SYSDATE)")
	case "A RECEBER", "ABERTO":
		conditions = append(conditions, "P.DTPAG IS NULL AND P.DTVENC >= TRUNC(SYSDATE)")
	}

	// Filtros por Data de Emissão (Apenas se informados)
	if strings.TrimSpace(f.DtEmissaoInicio) != "" {
		conditions = append(conditions, "P.DTEMISSAO >= TO_DATE(:dtemissao_inicio, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dtemissao_inicio", strings.TrimSpace(f.DtEmissaoInicio)))
	}
	if strings.TrimSpace(f.DtEmissaoFim) != "" {
		conditions = append(conditions, "P.DTEMISSAO <= TO_DATE(:dtemissao_fim, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dtemissao_fim", strings.TrimSpace(f.DtEmissaoFim)))
	}

	// Filtros por Data de Vencimento (Apenas se informados)
	if strings.TrimSpace(f.DtVencInicio) != "" {
		conditions = append(conditions, "P.DTVENC >= TO_DATE(:dtvenc_inicio, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dtvenc_inicio", strings.TrimSpace(f.DtVencInicio)))
	}
	if strings.TrimSpace(f.DtVencFim) != "" {
		conditions = append(conditions, "P.DTVENC <= TO_DATE(:dtvenc_fim, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dtvenc_fim", strings.TrimSpace(f.DtVencFim)))
	}

	// Filtros por Data de Pagamento / DTPAG (Apenas se informados)
	if strings.TrimSpace(f.DtPagInicio) != "" {
		conditions = append(conditions, "P.DTPAG >= TO_DATE(:dtpag_inicio, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dtpag_inicio", strings.TrimSpace(f.DtPagInicio)))
	}
	if strings.TrimSpace(f.DtPagFim) != "" {
		conditions = append(conditions, "P.DTPAG <= TO_DATE(:dtpag_fim, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dtpag_fim", strings.TrimSpace(f.DtPagFim)))
	}

	// Filtros por Faixa de Valor
	if f.ValorMin != nil {
		conditions = append(conditions, "P.VALOR >= :valor_min")
		args = append(args, sql.Named("valor_min", *f.ValorMin))
	}
	if f.ValorMax != nil {
		conditions = append(conditions, "P.VALOR <= :valor_max")
		args = append(args, sql.Named("valor_max", *f.ValorMax))
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " AND " + strings.Join(conditions, " AND ")
	}

	return whereClause, args
}

// ToDataSQL gera a instrução SQL com a lista de títulos, ordenação e limitação por ROWNUM
func (f FilterPCPrest) ToDataSQL() (string, []interface{}) {
	whereClause, args := f.BuildConditions()
	cte := fmt.Sprintf(QueryPCPrestCTE, whereClause)
	fullQuery := cte + "\n" + QueryPCPrestCalculada

	colunaOrdenacao := "DUPLIC, PREST"
	switch strings.ToLower(strings.TrimSpace(f.OrdenarPor)) {
	case "dtvenc", "vencimento":
		colunaOrdenacao = "DTVENC"
	case "dtemissao", "emissao":
		colunaOrdenacao = "DTEMISSAO"
	case "dtpag", "pagamento":
		colunaOrdenacao = "DTPAG"
	case "valor":
		colunaOrdenacao = "VALOR"
	case "cliente":
		colunaOrdenacao = "CLIENTE"
	}

	direcao := "ASC"
	if strings.ToUpper(strings.TrimSpace(f.Ordem)) == "DESC" {
		direcao = "DESC NULLS LAST"
	}

	limite := 50
	if f.Limite != nil && *f.Limite > 0 {
		limite = *f.Limite
	}

	fullQuery += fmt.Sprintf(" ORDER BY %s %s", colunaOrdenacao, direcao)
	fullQuery = fmt.Sprintf("SELECT * FROM (%s) WHERE ROWNUM <= %d", fullQuery, limite)

	return fullQuery, args
}

// ToCalculosSQL gera a instrução SQL para os Totais Agregados limitados ao mesmo ROWNUM do ToDataSQL
func (f FilterPCPrest) ToCalculosSQL() (string, []interface{}) {
	whereClause, args := f.BuildConditions()
	cte := fmt.Sprintf(QueryPCPrestCTE, whereClause)

	colunaOrdenacao := "DUPLIC, PREST"
	switch strings.ToLower(strings.TrimSpace(f.OrdenarPor)) {
	case "dtvenc", "vencimento":
		colunaOrdenacao = "DTVENC"
	case "dtemissao", "emissao":
		colunaOrdenacao = "DTEMISSAO"
	case "dtpag", "pagamento":
		colunaOrdenacao = "DTPAG"
	case "valor":
		colunaOrdenacao = "VALOR"
	case "cliente":
		colunaOrdenacao = "CLIENTE"
	}

	direcao := "ASC"
	if strings.ToUpper(strings.TrimSpace(f.Ordem)) == "DESC" {
		direcao = "DESC NULLS LAST"
	}

	limite := 50
	if f.Limite != nil && *f.Limite > 0 {
		limite = *f.Limite
	}

	summaryQuery := fmt.Sprintf(`%s
SELECT
    COUNT(*) AS QTD_TOTAL,
    COUNT(CASE WHEN DTPAG IS NOT NULL THEN 1 END) AS QTD_PAGO,
    COUNT(CASE WHEN DTPAG IS NULL AND DTVENC >= TRUNC(SYSDATE) THEN 1 END) AS QTD_ABERTO,
    COUNT(CASE WHEN DTPAG IS NULL AND DTVENC < TRUNC(SYSDATE) THEN 1 END) AS QTD_ATRASADO,
    NVL(SUM(VALOR), 0) AS VALOR_TOTAL,
    NVL(SUM(VPAGO), 0) AS VALOR_PAGO,
    NVL(SUM(CASE WHEN DTPAG IS NULL AND DTVENC >= TRUNC(SYSDATE) THEN VALOR ELSE 0 END), 0) AS VALOR_ABERTO,
    NVL(SUM(CASE WHEN DTPAG IS NULL AND DTVENC < TRUNC(SYSDATE) THEN VALOR ELSE 0 END), 0) AS VALOR_ATRASADO
FROM (
    SELECT D.* FROM DADOS D ORDER BY %s %s
) WHERE ROWNUM <= %d`, cte, colunaOrdenacao, direcao, limite)

	return summaryQuery, args
}

type PCPrestRepository interface {
	SearchPCPrest(c context.Context, filtro FilterPCPrest) ([]PCPrest, error)
	GetPCPrestCalculos(c context.Context, filtro FilterPCPrest) (*PCPrestSummary, error)
}

type PCPrestService interface {
	ListPCPrest(c context.Context, filtro FilterPCPrest) ([]PCPrest, *PCPrestSummary, error)
}
