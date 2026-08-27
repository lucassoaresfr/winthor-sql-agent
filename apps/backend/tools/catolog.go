package tools

import "google.golang.org/genai"

var ToolConsultarClientes = &genai.FunctionDeclaration{
	Name:        "consultar_clientes_api",
	Description: "Consulta a API de clientes/parceiros no WinThor permitindo filtrar por qualquer parâmetro cadastral, localização, rede, plano de pagamento, tipo de pessoa (PF/PJ) ou documento (CPF/CNPJ via cgcent). Utilize para verificar se um cliente, marca ou fornecedor está cadastrado.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			// --- Identificação e Filtros Principais ---
			"codcli": {
				Type:        genai.TypeInteger,
				Description: "Código único do cliente no WinThor (C.CODCLI).",
			},
			"cliente": {
				Type:        genai.TypeString,
				Description: "Razão social ou nome do cliente/fornecedor (C.CLIENTE).",
			},
			"fantasia": {
				Type:        genai.TypeString,
				Description: "Nome fantasia do cliente/marca/fornecedor (C.FANTASIA).",
			},
			"tipofj": {
				Type:        genai.TypeString,
				Description: "Tipo de pessoa no cadastro: 'F' para Física (CPF) ou 'J' para Jurídica (CNPJ) (C.TIPOFJ).",
			},
			"cgcent": {
				Type:        genai.TypeString,
				Description: "Documento de identificação cadastral do cliente/fornecedor (CPF ou CNPJ sem pontuação/com pontuação) (C.CGCENT).",
			},

			// --- Localização e Endereço ---
			"municent": {
				Type:        genai.TypeString,
				Description: "Município/cidade do cliente. Ex: RECIFE, OLINDA (C.MUNICENT).",
			},
			"bairroent": {
				Type:        genai.TypeString,
				Description: "Bairro do cliente. Ex: BOA VIAGEM, CENTRO (C.BAIRROENT).",
			},

			// --- Contato e Datas ---
			"telent": {
				Type:        genai.TypeString,
				Description: "Telefone do cliente (C.TELENT).",
			},
			"dtultcomp_inicio": {
				Type:        genai.TypeString,
				Description: "Data inicial da última compra no formato YYYY-MM-DD (C.DTULTCOMP).",
			},
			"dtultcomp_fim": {
				Type:        genai.TypeString,
				Description: "Data final da última compra no formato YYYY-MM-DD (C.DTULTCOMP).",
			},

			// --- Plano de Pagamento (PCPLPAG) ---
			"codplpag": {
				Type:        genai.TypeInteger,
				Description: "Código do plano de pagamento associado (P.CODPLPAG).",
			},
			"numpr": {
				Type:        genai.TypeInteger,
				Description: "Número de parcelas do plano de pagamento (P.NUMPR).",
			},
			"descricao_plpag": {
				Type:        genai.TypeString,
				Description: "Descrição do plano de pagamento (P.DESCRICAO).",
			},
			"numdias": {
				Type:        genai.TypeInteger,
				Description: "Prazo de dias do plano de pagamento (P.NUMDIAS).",
			},

			// --- Rede de Clientes (PCREDECLIENTE) ---
			"codrede": {
				Type:        genai.TypeInteger,
				Description: "Código da rede do cliente (C.CODREDE).",
			},
			"descricao_rede": {
				Type:        genai.TypeString,
				Description: "Nome/Descrição da rede de clientes (R.DESCRICAO).",
			},

			// --- Outros Filtros WinThor ---
			"codcob": {
				Type:        genai.TypeString,
				Description: "Código da cobrança cadastrada. Ex: D, BK, CH (C.CODCOB).",
			},

			// --- Ordenação / Tipo de ordenação
			"ordenar_por": {
				Type:        genai.TypeString,
				Description: "Campo pelo qual deseja ordenar os resultados. Opções permitidas: 'dtultcomp', 'cliente', 'fantasia', 'codcli', 'municent'. Padrão: 'codcli'.",
			},
			"ordem": {
				Type:        genai.TypeString,
				Description: "Direção da ordenação: 'ASC' ou 'DESC'. Padrão: 'ASC'.",
			},
			"limite": {
				Type:        genai.TypeInteger,
				Description: "Quantidade máxima de registros a retornar. Padrão: 20.",
			},
		},
	},
}

var ToolConsultarProdutos = &genai.FunctionDeclaration{
	Name:        "consultar_produtos_api",
	Description: "Consulta o catálogo de produtos no ERP WinThor permitindo buscar por código, descrição, embalagem, unidade, fornecedor, departamento, marca e saldo de estoque da Filial 1.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			// --- Identificação do Produto ---
			"codprod": {
				Type:        genai.TypeInteger,
				Description: "Código único do produto no WinThor (P.CODPROD).",
			},
			"descricaoprod": {
				Type:        genai.TypeString,
				Description: "Descrição ou nome do produto (P.DESCRICAO). Aceita buscas parciais.",
			},
			"embalagem": {
				Type:        genai.TypeString,
				Description: "Tipo de embalagem do produto (P.EMBALAGEM), ex: 'CX', 'UN', 'FD', 'KG'. Aceita buscas parciais.",
			},
			"unidade": {
				Type:        genai.TypeString,
				Description: "Unidade de medida comercial (P.UNIDADE), ex: 'UN', 'CX', 'KG', 'L'. Aceita buscas parciais.",
			},

			// --- Fornecedor (PCFORNEC) ---
			"codfornec": {
				Type:        genai.TypeInteger,
				Description: "Código do fornecedor do produto (P.CODFORNEC).",
			},
			"fornecedor": {
				Type:        genai.TypeString,
				Description: "Nome/Razão Social do fornecedor (F.FORNECEDOR). Aceita buscas parciais.",
			},

			// --- Departamento (PCDEPTO) ---
			"coddepto": {
				Type:        genai.TypeInteger,
				Description: "Código do departamento do produto (P.CODEPTO).",
			},
			"descricao_depto": {
				Type:        genai.TypeString,
				Description: "Nome ou descrição do departamento: CONGELADOS, SECOS, RESFRIADOS (D.DESCRICAO). Aceita buscas parciais.",
			},

			// --- Marca (PCMARCA) ---
			"codmarca": {
				Type:        genai.TypeInteger,
				Description: "Código da marca do produto (P.CODMARCA).",
			},
			"marca": {
				Type:        genai.TypeString,
				Description: "Nome da marca do produto (M.MARCA). Aceita buscas parciais.",
			},

			// --- Filtro de Estoque ---
			"apenas_estoque": {
				Type:        genai.TypeBoolean,
				Description: "Se 'true', retorna apenas produtos com saldo disponível de estoque maior que zero (ESTOQUE > 0) na Filial 1.",
			},

			// --- Controle de Ordenação e Paginação ---
			"ordenar_por": {
				Type:        genai.TypeString,
				Description: "Campo para ordenação dos resultados. Opções permitidas: 'descricaoprod' (nome do produto), 'estoque' (quantidade de estoque) ou 'codprod' (código). Padrão: 'codprod'.",
			},
			"ordem": {
				Type:        genai.TypeString,
				Description: "Direção da ordenação: 'ASC' (crescente/alfabética) ou 'DESC' (decrescente/maior valor). Use 'DESC' para consultas sobre 'maior estoque'. Padrão: 'ASC'.",
			},
			"limite": {
				Type:        genai.TypeInteger,
				Description: "Quantidade máxima de registros a retornar (ex: 5, 10, 50). Padrão: 50.",
			},
		},
	},
}

var ToolConsultarCNPJExterno = &genai.FunctionDeclaration{
	Name:        "consultar_cnpj_externo",
	Description: "Consulta dados cadastrais públicos de uma empresa pelo CNPJ na Receita Federal via BrasilAPI. Use esta ferramenta quando o usuário quiser informações sobre uma empresa/marca antes de cadastrá-la no WinThor ou para validar Razão Social, Nome Fantasia, CNAE e Endereço público.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"cnpj": {
				Type:        genai.TypeString,
				Description: "O número do CNPJ contendo 14 dígitos (apenas números ou formatado ex: '00.000.000/0001-00').",
			},
		},
		Required: []string{"cnpj"},
	},
}
