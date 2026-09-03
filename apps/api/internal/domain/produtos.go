package domain

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Produto struct {
	CodProd        int     `json:"codprod" gorm:"column:CODPROD"`
	DescricaoProd  string  `json:"descricaoprod" gorm:"column:DESCRICAOPROD"`
	Embalagem      *string `json:"embalagem" gorm:"column:EMBALAGEM"`
	Unidade        *string `json:"unidade" gorm:"column:UNIDADE"`
	CodFornec      *int    `json:"codfornec" gorm:"column:CODFORNEC"`
	Fornecedor     *string `json:"fornecedor" gorm:"column:FORNECEDOR"`
	CodDepto       *int    `json:"coddepto" gorm:"column:CODEPTO"`
	DescricaoDepto *string `json:"descricao_depto" gorm:"column:DESCRICAO_DEPTO"`
	CodSecao       *int    `json:"codsecao" gorm:"column:CODSECAO"`
	DescricaoSecao *string `json:"descricao_secao" gorm:"column:DESCRICAOSECAO"`
	CodMarca       *int    `json:"codmarca" gorm:"column:CODMARCA"`
	Marca          *string `json:"marca" gorm:"column:MARCA"`
	Estoque        float64 `json:"estoque" gorm:"column:ESTOQUE"`
}

type FiltroProduto struct {
	CodProd        *int   `form:"codprod"`
	DescricaoProd  string `form:"descricaoprod"`
	Embalagem      string `form:"embalagem"`
	Unidade        string `form:"unidade"`
	CodFornec      *int   `form:"codfornec"`
	Fornecedor     string `form:"fornecedor"`
	CodDepto       *int   `form:"coddepto"`
	DescricaoDepto string `form:"descricao_depto"`
	CodSecao       *int   `form:"codsecao"`
	DescricaoSecao string `form:"descricao_secao"`
	CodMarca       *int   `form:"codmarca"`
	Marca          string `form:"marca"`
	ApenasEstoque  *bool  `form:"apenas_estoque"`
	OrdenarPor     string `form:"ordenar_por"`
	Ordem          string `form:"ordem"`
	Limite         *int   `form:"limite"`
}

const QueryProdutoBase = `
SELECT 
    P.CODPROD, 
    P.DESCRICAO AS DESCRICAOPROD, 
    P.EMBALAGEM, 
    P.UNIDADE, 
    P.CODFORNEC, 
    F.FORNECEDOR, 
    P.CODEPTO, 
    D.DESCRICAO AS DESCRICAO_DEPTO,
    P.CODSEC,
    S.DESCRICAO AS DESCRICAOSECAO,
    P.CODMARCA,
    M.MARCA,
    (E.QTRESERV - (E.QTBLOQUEADA - E.QTESTGER)) AS ESTOQUE
FROM PCPRODUT P
	LEFT JOIN PCFORNEC F ON P.CODFORNEC = F.CODFORNEC
	LEFT JOIN PCDEPTO D  ON P.CODEPTO = D.CODEPTO
	LEFT JOIN PCMARCA M  ON P.CODMARCA = M.CODMARCA
  LEFT JOIN PCSECAO S ON P.CODSEC = S.CODSEC
	LEFT JOIN PCEST E ON P.CODPROD = E.CODPROD AND E.CODFILIAL = '1'
WHERE P.DTEXCLUSAO IS NULL`

// ToSQL constrói a instrução SQL com parâmetros nomeados compatíveis com Oracle
func (f FiltroProduto) ToSQL() (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if f.CodProd != nil {
		conditions = append(conditions, "P.CODPROD = :codprod")
		args = append(args, sql.Named("codprod", *f.CodProd))
	}
	if strings.TrimSpace(f.DescricaoProd) != "" {
		conditions = append(conditions, "UPPER(P.DESCRICAO) LIKE UPPER(:descricaoprod)")
		args = append(args, sql.Named("descricaoprod", "%"+strings.TrimSpace(f.DescricaoProd)+"%"))
	}
	if strings.TrimSpace(f.Embalagem) != "" {
		conditions = append(conditions, "UPPER(P.EMBALAGEM) LIKE UPPER(:embalagem)")
		args = append(args, sql.Named("embalagem", "%"+strings.TrimSpace(f.Embalagem)+"%"))
	}
	if strings.TrimSpace(f.Unidade) != "" {
		conditions = append(conditions, "UPPER(P.UNIDADE) LIKE UPPER(:unidade)")
		args = append(args, sql.Named("unidade", "%"+strings.TrimSpace(f.Unidade)+"%"))
	}
	if f.CodFornec != nil {
		conditions = append(conditions, "P.CODFORNEC = :codfornec")
		args = append(args, sql.Named("codfornec", *f.CodFornec))
	}
	if strings.TrimSpace(f.Fornecedor) != "" {
		conditions = append(conditions, "UPPER(F.FORNECEDOR) LIKE UPPER(:fornecedor)")
		args = append(args, sql.Named("fornecedor", "%"+strings.TrimSpace(f.Fornecedor)+"%"))
	}
	if f.CodDepto != nil {
		conditions = append(conditions, "P.CODEPTO = :coddepto")
		args = append(args, sql.Named("coddepto", *f.CodDepto))
	}
	if strings.TrimSpace(f.DescricaoDepto) != "" {
		conditions = append(conditions, "UPPER(D.DESCRICAO) LIKE UPPER(:descricao_depto)")
		args = append(args, sql.Named("descricao_depto", "%"+strings.TrimSpace(f.DescricaoDepto)+"%"))
	}
	if f.CodSecao != nil {
		conditions = append(conditions, "P.CODSECAO = :codsecao")
		args = append(args, sql.Named("codsecao", *f.CodSecao))
	}
	if strings.TrimSpace(f.DescricaoSecao) != "" {
		conditions = append(conditions, "UPPER(S.DESCRICAO) LIKE UPPER(:descricao_secao)")
		args = append(args, sql.Named("descricao_secao", "%"+strings.TrimSpace(f.DescricaoSecao)+"%"))
	}
	if f.CodMarca != nil {
		conditions = append(conditions, "P.CODMARCA = :codmarca")
		args = append(args, sql.Named("codmarca", *f.CodMarca))
	}
	if strings.TrimSpace(f.Marca) != "" {
		conditions = append(conditions, "UPPER(M.MARCA) LIKE UPPER(:marca)")
		args = append(args, sql.Named("marca", "%"+strings.TrimSpace(f.Marca)+"%"))
	}
	if f.ApenasEstoque != nil && *f.ApenasEstoque {
		conditions = append(conditions, "(NVL(E.QTESTGER, 0) - NVL(E.QTRESERV, 0) - NVL(E.QTBLOQUEADA, 0)) > 0")
	}

	query := QueryProdutoBase
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	colunaOrdenacao := "P.CODPROD"
	switch strings.ToLower(strings.TrimSpace(f.OrdenarPor)) {
	case "descricaoprod", "descricao":
		colunaOrdenacao = "P.DESCRICAO"
	case "estoque":
		colunaOrdenacao = "ESTOQUE"
	case "codprod":
		colunaOrdenacao = "P.CODPROD"
	}

	direcao := "ASC"
	if strings.ToUpper(strings.TrimSpace(f.Ordem)) == "DESC" {
		direcao = "DESC NULLS LAST"
	}

	limite := 50
	if f.Limite != nil && *f.Limite > 0 {
		limite = *f.Limite
	}

	query += " ORDER BY " + colunaOrdenacao + " " + direcao
	query = fmt.Sprintf("SELECT * FROM (%s) WHERE ROWNUM <= %d", query, limite)

	return query, args
}

type ProdutoRepository interface {
	BuscarProdutos(c context.Context, filtro FiltroProduto) ([]Produto, error)
}

type ProdutoService interface {
	ListarProdutos(c context.Context, filtro FiltroProduto) ([]Produto, error)
}
