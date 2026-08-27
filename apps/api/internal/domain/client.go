package domain

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/pkg"
)

type Cliente struct {
	CodCli         int      `json:"codcli" gorm:"column:CODCLI"`
	Cliente        string   `json:"cliente" gorm:"column:CLIENTE"`
	Fantasia       *string  `json:"fantasia" gorm:"column:FANTASIA"`
	TelEnt         *string  `json:"telent" gorm:"column:TELENT"`
	DtUltComp      *string  `json:"dtultcomp" gorm:"column:DTULTCOMP"`
	BairroEnt      *string  `json:"bairroent" gorm:"column:BAIRROENT"`
	TipoFJ         *string  `json:"tipofj" gorm:"column:TIPOFJ"`
	MunicEnt       *string  `json:"municent" gorm:"column:MUNICENT"`
	CodPlPag       *int     `json:"codplpag" gorm:"column:CODPLPAG"`
	NumPR          *int     `json:"numpr" gorm:"column:NUMPR"`
	DescricaoPlPag *string  `json:"descricao_plpag" gorm:"column:DESCRICAO_PLPAG"`
	CodRede        *int     `json:"codrede" gorm:"column:CODREDE"`
	DescricaoRede  *string  `json:"descricao_rede" gorm:"column:DESCRICAO_REDE"`
	PerTxFim       *float64 `json:"pertxfim" gorm:"column:PERTXFIM"`
	NumDias        *int     `json:"numdias" gorm:"column:NUMDIAS"`
	CodCob         *string  `json:"codcob" gorm:"column:CODCOB"`
	Cgcent         *string  `json:"cgcent" gorm:"column:CGCENT"`
}

// MascararDadosSensiveis aplica a política LGPD de mascaramento para Pessoas Físicas (TIPOFJ == "F")
func (c *Cliente) MascararDadosSensiveis() {
	if c.TipoFJ != nil && strings.ToUpper(strings.TrimSpace(*c.TipoFJ)) == "F" {
		// 1. Mascara o Nome Razão Social
		if c.Cliente != "" {
			c.Cliente = pkg.MascararNomePessoaFisica(c.Cliente)
		}

		// 2. Mascara Nome Fantasia
		if c.Fantasia != nil && *c.Fantasia != "" {
			fantasiaOculta := pkg.MascararNomePessoaFisica(*c.Fantasia)
			c.Fantasia = &fantasiaOculta
		}

		// 3. Mascara o Telefone
		if c.TelEnt != nil && *c.TelEnt != "" {
			telOculto := pkg.MascararTelefone(*c.TelEnt)
			c.TelEnt = &telOculto
		}

		// 4. Mascara o Documento (CPF)
		if c.Cgcent != nil && *c.Cgcent != "" {
			docOculto := pkg.MascararDocumento(*c.Cgcent)
			c.Cgcent = &docOculto
		}
	}
}

type FiltroCliente struct {
	CodCli          *int    `form:"codcli"`
	Cliente         string  `form:"cliente"`
	Fantasia        string  `form:"fantasia"`
	TipoFJ          string  `form:"tipofj"`
	MunicEnt        string  `form:"municent"`
	BairroEnt       string  `form:"bairroent"`
	TelEnt          string  `form:"telent"`
	DtUltCompInicio string  `form:"dtultcomp_inicio"`
	DtUltCompFim    string  `form:"dtultcomp_fim"`
	CodPlPag        *int    `form:"codplpag"`
	NumPR           *int    `form:"numpr"`
	NumDias         *int    `form:"numdias"`
	CodRede         *int    `form:"codrede"`
	CodCob          *string `form:"codcob"`
	Cgcent          *string `form:"cgcent"`
	OrdenarPor      string  `form:"ordenar_por"`
	Ordem           string  `form:"ordem"`
	Limite          *int    `form:"limite"`
}

const QueryClientBase = `SELECT C.CODCLI, C.CLIENTE, C.FANTASIA, 
C.TELENT, TO_CHAR(C.DTULTCOMP, 'YYYY-MM-DD') AS DTULTCOMP, C.BAIRROENT, 
C.TIPOFJ, C.MUNICENT, P.CODPLPAG, P.NUMPR, P.DESCRICAO AS DESCRICAO_PLPAG, 
C.CODREDE, R.DESCRICAO AS DESCRICAO_REDE, P.PERTXFIM, P.NUMDIAS, C.CODCOB, C.CGCENT 
FROM PCCLIENT C 
LEFT JOIN PCPLPAG P ON C.CODPLPAG = P.CODPLPAG 
LEFT JOIN PCREDECLIENTE R ON C.CODREDE = R.CODREDE WHERE 1=1`

// ToSQL constrói a instrução SQL com parâmetros nomeados compatíveis com Oracle
func (f FiltroCliente) ToSQL() (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if f.CodCli != nil {
		conditions = append(conditions, "C.CODCLI = :codcli")
		args = append(args, sql.Named("codcli", *f.CodCli))
	}
	if strings.TrimSpace(f.Cliente) != "" {
		conditions = append(conditions, "UPPER(C.CLIENTE) LIKE UPPER(:cliente)")
		args = append(args, sql.Named("cliente", "%"+strings.TrimSpace(f.Cliente)+"%"))
	}
	if strings.TrimSpace(f.Fantasia) != "" {
		conditions = append(conditions, "UPPER(C.FANTASIA) LIKE UPPER(:fantasia)")
		args = append(args, sql.Named("fantasia", "%"+strings.TrimSpace(f.Fantasia)+"%"))
	}
	if strings.TrimSpace(f.TipoFJ) != "" {
		conditions = append(conditions, "UPPER(C.TIPOFJ) = UPPER(:tipofj)")
		args = append(args, sql.Named("tipofj", strings.TrimSpace(f.TipoFJ)))
	}
	if strings.TrimSpace(f.MunicEnt) != "" {
		conditions = append(conditions, "UPPER(C.MUNICENT) LIKE UPPER(:municent)")
		args = append(args, sql.Named("municent", "%"+strings.TrimSpace(f.MunicEnt)+"%"))
	}
	if strings.TrimSpace(f.BairroEnt) != "" {
		conditions = append(conditions, "UPPER(C.BAIRROENT) LIKE UPPER(:bairroent)")
		args = append(args, sql.Named("bairroent", "%"+strings.TrimSpace(f.BairroEnt)+"%"))
	}
	if strings.TrimSpace(f.TelEnt) != "" {
		conditions = append(conditions, "C.TELENT LIKE :telent")
		args = append(args, sql.Named("telent", "%"+strings.TrimSpace(f.TelEnt)+"%"))
	}
	if strings.TrimSpace(f.DtUltCompInicio) != "" {
		conditions = append(conditions, "C.DTULTCOMP >= TO_DATE(:dt_inicio, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dt_inicio", strings.TrimSpace(f.DtUltCompInicio)))
	}
	if strings.TrimSpace(f.DtUltCompFim) != "" {
		conditions = append(conditions, "C.DTULTCOMP <= TO_DATE(:dt_fim, 'YYYY-MM-DD')")
		args = append(args, sql.Named("dt_fim", strings.TrimSpace(f.DtUltCompFim)))
	}
	if f.CodPlPag != nil {
		conditions = append(conditions, "P.CODPLPAG = :codplpag")
		args = append(args, sql.Named("codplpag", *f.CodPlPag))
	}
	if f.NumPR != nil {
		conditions = append(conditions, "P.NUMPR = :numpr")
		args = append(args, sql.Named("numpr", *f.NumPR))
	}
	if f.NumDias != nil {
		conditions = append(conditions, "P.NUMDIAS = :numdias")
		args = append(args, sql.Named("numdias", *f.NumDias))
	}
	if f.CodRede != nil {
		conditions = append(conditions, "C.CODREDE = :codrede")
		args = append(args, sql.Named("codrede", *f.CodRede))
	}
	if f.CodCob != nil && strings.TrimSpace(*f.CodCob) != "" {
		conditions = append(conditions, "UPPER(C.CODCOB) = UPPER(:codcob)")
		args = append(args, sql.Named("codcob", strings.TrimSpace(*f.CodCob)))
	}
	if f.Cgcent != nil && strings.TrimSpace(*f.Cgcent) != "" {
		conditions = append(conditions, "UPPER(C.CGCENT) = UPPER(:cgcent)")
		args = append(args, sql.Named("cgcent", strings.TrimSpace(*f.Cgcent)))
	}

	query := QueryClientBase
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	colunaOrdenacao := "C.CODCLI"
	switch strings.ToLower(strings.TrimSpace(f.OrdenarPor)) {
	case "dtultcomp":
		colunaOrdenacao = "C.DTULTCOMP"
	case "cliente":
		colunaOrdenacao = "C.CLIENTE"
	case "fantasia":
		colunaOrdenacao = "C.FANTASIA"
	case "municent":
		colunaOrdenacao = "C.MUNICENT"
	case "codcli":
		colunaOrdenacao = "C.CODCLI"
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

type ClientRepository interface {
	BuscarClientes(ctx context.Context, filtro FiltroCliente) ([]Cliente, error)
}

type ClientService interface {
	ListarClientes(ctx context.Context, filtro FiltroCliente) ([]Cliente, error)
}
