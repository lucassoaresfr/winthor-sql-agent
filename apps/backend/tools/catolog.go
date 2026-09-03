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

var ToolConsultarPedidos = &genai.FunctionDeclaration{
	Name:        "consultar_pedidos_api",
	Description: "Consulta a capa dos pedidos de venda no WinThor (PCPEDC/Orders). Permite filtrar por número do pedido, cliente, vendedor/RCA, posição/status, valores, plano de pagamento, cobrança, condição de venda, origem, carregamento, data de faturamento, observações e período de emissão. Com 'incluir_itens=true', retorna também a lista de itens/produtos (ItemOrder) aninhada em cada pedido.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			// --- Identificação do Pedido ---
			"numped": {
				Type:        genai.TypeInteger,
				Description: "Número do pedido de venda no WinThor (PCPEDC.NUMPED).",
			},

			// --- Cliente ---
			"codcli": {
				Type:        genai.TypeInteger,
				Description: "Código único do cliente (PCPEDC.CODCLI).",
			},
			"cliente": {
				Type:        genai.TypeString,
				Description: "Razão social ou parte do nome do cliente (PCPEDC.CLIENTE / PCCLIENT.CLIENTE).",
			},

			// --- Vendedor / RCA ---
			"codusur": {
				Type:        genai.TypeInteger,
				Description: "Código do vendedor/RCA (PCPEDC.CODUSUR).",
			},
			"nomevend": {
				Type:        genai.TypeString,
				Description: "Nome ou parte do nome do vendedor/RCA (PCUSUARI.NOME).",
			},

			// --- Status e Condições Comerciais ---
			"posicao": {
				Type:        genai.TypeString,
				Description: "Posição/status do pedido no WinThor (PCPEDC.POSICAO). Ex: 'F' (Faturado), 'L' (Liberado), 'P' (Pendente), 'C' (Cancelado), 'M' (Montado).",
			},
			"codplpag": {
				Type:        genai.TypeInteger,
				Description: "Código do plano de pagamento (PCPEDC.CODPLPAG).",
			},
			"descricaopag": {
				Type:        genai.TypeString,
				Description: "Descrição do plano de pagamento associado ao pedido.",
			},
			"codcob": {
				Type:        genai.TypeString,
				Description: "Código da forma de cobrança/pagamento. Ex: 'D' (Dinheiro), 'BK' (Bancário), 'CH' (Cheque) (PCPEDC.CODCOB).",
			},
			"condvenda": {
				Type:        genai.TypeInteger,
				Description: "Condição de venda no WinThor (PCPEDC.CONDVENDA). Ex: 1 para venda normal, 5 para bonificação.",
			},
			"origemped": {
				Type:        genai.TypeString,
				Description: "Origem da emissão do pedido (PCPEDC.ORIGEMPED). Ex: 'F' (Força de vendas), 'W' (Web), 'T' (Balcão).",
			},

			// --- Logística e Faturamento ---
			"numcar": {
				Type:        genai.TypeInteger,
				Description: "Número do carregamento de entrega do pedido (PCPEDC.NUMCAR).",
			},
			"dtfat": {
				Type:        genai.TypeString,
				Description: "Data de faturamento do pedido no formato YYYY-MM-DD (PCPEDC.DTFAT).",
			},
			"obs": {
				Type:        genai.TypeString,
				Description: "Texto de observação do pedido (PCPEDC.OBS).",
			},

			// --- Filtros de Período de Emissão ---
			"dt_inicio": {
				Type:        genai.TypeString,
				Description: "Data inicial de emissão do pedido no formato YYYY-MM-DD (PCPEDC.DATA >= dt_inicio).",
			},
			"dt_fim": {
				Type:        genai.TypeString,
				Description: "Data final de emissão do pedido no formato YYYY-MM-DD (PCPEDC.DATA <= dt_fim).",
			},

			// --- Inclusão de Itens Aninhados ---
			"incluir_itens": {
				Type:        genai.TypeBoolean,
				Description: "Quando true, a API realiza a busca complementar e popula a chave 'itens' com a lista de ItemOrder (numped, codprod, descricao, qt, pvenda, ptabela, perdesc) de cada pedido.",
			},

			// --- Ordenação e Limites ---
			"ordenar_por": {
				Type:        genai.TypeString,
				Description: "Campo para ordenação dos pedidos. Opções permitidas: 'data', 'numped', 'vltotal', 'cliente'. Padrão: 'numped'.",
			},
			"ordem": {
				Type:        genai.TypeString,
				Description: "Direção da ordenação: 'ASC' ou 'DESC'. Padrão: 'DESC'.",
			},
			"limite": {
				Type:        genai.TypeInteger,
				Description: "Quantidade máxima de pedidos a retornar. Padrão: 20.",
			},
		},
	},
}

var ToolConsultarItensPedido = &genai.FunctionDeclaration{
	Name:        "consultar_itens_pedido_api",
	Description: "Consulta os itens e produtos pertencentes aos pedidos de venda no WinThor (PCPEDI/ItemOrder). Permite filtrar por número do pedido (numped), código do produto (codprod), descrição do produto, seção (código ou descrição), quantidade vendida (qt), preços de venda (pvenda), preço de tabela (ptabela) e percentual de desconto (perdesc).",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			// --- Identificação do Pedido e Produto ---
			"numped": {
				Type:        genai.TypeInteger,
				Description: "Número do pedido de venda para listar seus itens (PCPEDI.NUMPED).",
			},
			"codprod": {
				Type:        genai.TypeInteger,
				Description: "Código do produto no WinThor (PCPEDI.CODPROD). Use para verificar em quais pedidos este produto foi vendido.",
			},
			"descricao": {
				Type:        genai.TypeString,
				Description: "Descrição ou nome parcial do produto (PCPRODUT.DESCRICAO).",
			},

			// --- Seção do Produto ---
			"codsecao": {
				Type:        genai.TypeInteger,
				Description: "Código da seção do produto (PCPRODUT.CODSECAO). Use para listar vendas de uma seção específica.",
			},
			"descricao_secao": {
				Type:        genai.TypeString,
				Description: "Descrição ou nome parcial da seção do produto (PCSECAO.DESCRICAO). Ex: 'BEBIDAS', 'HORTIFRUTI', 'FROZEN'.",
			},

			// --- Quantidade ---
			"qt": {
				Type:        genai.TypeNumber,
				Description: "Quantidade exata ou mínima vendida do produto no item do pedido (PCPEDI.QT).",
			},

			// --- Preços e Valores ---
			"pvenda": {
				Type:        genai.TypeNumber,
				Description: "Preço unitário praticado de venda do produto no item (PCPEDI.PVENDA).",
			},
			"ptabela": {
				Type:        genai.TypeNumber,
				Description: "Preço unitário de tabela do produto no item (PCPEDI.PTABELA).",
			},
			"perdesc": {
				Type:        genai.TypeNumber,
				Description: "Percentual de desconto aplicado no item do pedido (PCPEDI.PERDESC). Exemplo: 5.0 para 5% de desconto.",
			},

			// --- Ordenação e Limites ---
			"ordenar_por": {
				Type:        genai.TypeString,
				Description: "Campo para ordenação dos itens. Opções permitidas: 'numped', 'codprod', 'qt', 'pvenda', 'perdesc'. Padrão: 'numped'.",
			},
			"ordem": {
				Type:        genai.TypeString,
				Description: "Direção da ordenação: 'ASC' ou 'DESC'. Padrão: 'ASC'.",
			},
			"limite": {
				Type:        genai.TypeInteger,
				Description: "Quantidade máxima de itens a retornar. Padrão: 50.",
			},
		},
	},
}

var ToolConsultarPromocoes = &genai.FunctionDeclaration{
	Name:        "consultar_promocoes_api",
	Description: "Consulta as promoções ativas no sistema WinThor (PCPRECOPROM/PCPRODUT/PCEST). Permite listar produtos em promoção com busca parcial, filtrar por estoque disponível via 'apenas_estoque: true', ou agrupar apenas as modalidades/descrições de ofertas ativas utilizando 'apenas_tipos: true'.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			// --- Identificadores ---
			"codprecoprom": {
				Type:        genai.TypeInteger,
				Description: "Código identificador exato da promoção de preço no WinThor (PCPRECOPROM.CODPRECOPROM).",
			},
			"coddescricao": {
				Type:        genai.TypeInteger,
				Description: "Código do agrupador ou descrição da promoção (PCPRECOPROM.CODDESCRICAO).",
			},
			"codprod": {
				Type:        genai.TypeInteger,
				Description: "Código do produto em promoção (PCPRECOPROM.CODPROD).",
			},

			// --- Buscas por Texto (Qualquer parte da string) ---
			"descpromo": {
				Type:        genai.TypeString,
				Description: "Nome ou trecho da descrição da promoção (PCPRECOPROM.DESCRICAO). Ex: 'OFERTA', 'BLACK', 'DESCONTO'. Busca parcial.",
			},
			"descprod": {
				Type:        genai.TypeString,
				Description: "Descrição ou nome parcial do produto (PCPRODUT.DESCRICAO). Ex: 'FEIJAO', 'LEITE', 'DETERGENTE'. Busca parcial.",
			},
			"embalagem": {
				Type:        genai.TypeString,
				Description: "Tipo de embalagem do produto (PCPRODUT.EMBALAGEM). Ex: 'CX', 'FARDO', 'UN', 'PCT'. Busca parcial.",
			},
			"unidade": {
				Type:        genai.TypeString,
				Description: "Unidade de medida do produto (PCPRODUT.UNIDADE). Ex: 'KG', 'UN', 'LT'. Busca parcial.",
			},

			// --- Preço, Estoque e Modalidades ---
			"precofixo": {
				Type:        genai.TypeNumber,
				Description: "Valor exato do preço promocional fixado (PCPRECOPROM.PRECOFIXO). Ex: 10.50.",
			},
			"apenas_estoque": {
				Type:        genai.TypeBoolean,
				Description: "Se 'true', filtra apenas promoções com estoque gerencial disponível maior que zero (QTESTGER - QTRESERV - QTBLOQUEADA > 0).",
			},
			"apenas_tipos": {
				Type:        genai.TypeBoolean,
				Description: "Se 'true', retorna apenas a lista distinta (DISTINCT) das modalidades/descrições de ofertas ativas no ERP (ex: 'PROMOCAO DO DIA', 'QUEIMA DE ESTOQUE'), sem detalhar produtos individuais. Ideal para responder perguntas como 'quais tipos de ofertas temos hoje?'.",
			},

			// --- Ordenação e Limite ---
			"ordenar_por": {
				Type:        genai.TypeString,
				Description: "Campo para ordenação. Opções: 'codprecoprom', 'codprod', 'descprod', 'descpromo', 'precofixo', 'estoque', 'embalagem', 'unidade'. Padrão: 'codprecoprom'.",
			},
			"ordem": {
				Type:        genai.TypeString,
				Description: "Direção da ordenação: 'ASC' ou 'DESC'. Padrão: 'DESC'.",
			},
			"limite": {
				Type:        genai.TypeInteger,
				Description: "Quantidade máxima de registros retornados. Padrão: 100.",
			},
		},
	},
}
