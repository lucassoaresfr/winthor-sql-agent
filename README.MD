<div align="center">

  # 🤖 WinThor Agent AI

  <p align="center">
    <b>Assistente Virtual Inteligente integrado ao ERP WinThor & PostgreSQL</b>
  </p>

  <!-- Sugestão de Imagem 1: Banner Principal / Demo Animada -->
  <a href="#-demonstração">
    <img src="docs/banner-preview.png" alt="WinThor Agent AI Banner" width="100%" style="border-radius: 8px;">
  </a>

  <p align="center">
    <img src="https://img.shields.io/badge/Next.js-15-black?style=for-the-badge&logo=next.js" alt="Next.js" />
    <img src="https://img.shields.io/badge/React-19-blue?style=for-the-badge&logo=react" alt="React" />
    <img src="https://img.shields.io/badge/TypeScript-5.0-blue?style=for-the-badge&logo=typescript" alt="TypeScript" />
    <img src="https://img.shields.io/badge/Tailwind_CSS-3.4-38B2AC?style=for-the-badge&logo=tailwind-css" alt="Tailwind" />
  </p>

</div>

---

## 📌 Sobre o Projeto

O **WinThor Agent AI** é um assistente virtual corporativo desenvolvido para otimizar o fluxo de trabalho de equipes que utilizam o sistema ERP WinThor. 

A plataforma permite que usuários façam consultas em tempo real utilizando linguagem natural para obter informações cruciais sobre produtos, saldos de estoque, pedidos de venda, dados cadastrais de clientes e ofertas ativas, unindo a inteligência artificial da Google (Gemini) ao banco de dados do sistema.

### ✨ Principais Funcionalidades

- **📦 Produtos & Estoque:** Consulta instantânea de preços, saldos disponíveis para venda e categorização (*Departamento / Seção*).
- **🛒 Pedidos de Venda:** Verificação de status de pedidos (*Faturado, Liberado, Pendente*), resumo de itens e valores.
- **🏷️ Promoções & Ofertas:** Consulta de descontos vigentes e campanhas ativas.
- **🏢 Clientes & CNPJ:** Busca rápida de dados cadastrais internos e validação pública via Receita Federal.
- **💬 Histórico Persistente:** Suporte a conversas contínuas com contexto de histórico gravado em PostgreSQL via backend em Go/GORM.

---

## 📸 Demonstração Visal

<!-- Sugestão de Imagem 2: GIF do fluxo de uso -->
<div align="center">
  <img src="docs/demo-flow.gif" alt="Demonstração em Vídeo/GIF" width="85%" style="border-radius: 8px; border: 1px solid #e2e8f0;">
</div>

---

## 🛠️ Métodos e Arquitetura do Sistema

A aplicação combina **Next.js (App Router)** no front-end e um serviço **REST API** para persistência do histórico e integração com o modelo de linguagem.

### 1. Serviço de Histórico (`chatService`)
Módulo responsável por gerenciar sessões e mensagens no banco de dados.

| Método | Tipo | Rota / Função | Descrição |
| :--- | :--- | :--- | :--- |
| `listUserChats` | `GET` | `/chats/user/:userId` | Recupera todas as conversas salvas associadas a um usuário. |
| `getChatById` | `GET` | `/chats/:id` | Carrega o histórico de mensagens de uma conversa específica. |
| `createChat` | `POST` | `/chats` | Cria uma nova sessão de chat com o ID e o título inicial. |
| `addMessage` | `POST` | `/chats/:id/messages` | Adiciona uma nova mensagem (`user` ou `assistant`) ao histórico. |

### 2. Serviço de Inteligência (`ChatApi`)
Módulo responsável por se comunicar com a IA e retornar a resposta formatada.

| Método | Parâmetros | Descrição |
| :--- | :--- | :--- |
| `ChatApi` | `{ messages: Message[] }` | Recebe a janela de contexto atual (máximo de 6 mensagens) e retorna a resposta gerada pelo modelo Gemini. |

---

## 💡 Onde Colocar Imagens no Repositório (Locais Sugeridos)

Para garantir que seu README seja visualmente atraente e organizado, adicione as imagens na pasta `public/docs/` ou `docs/` e insira-as nos seguintes pontos estratégicos:

### 1. Banner Principal (No topo do README)
- **Onde:** Logo abaixo do título e badges.
- **O que colocar:** Uma captura da tela cheia do sistema aberto em uma conversa bonita ou um banner estilizado no Figma/Canva com a logo do WinThor Agent AI.

### 2. Demonstração em Ação (GIF ou Vídeo Curto)
- **Onde:** Na seção `Demonstração Visual`.
- **O que colocar:** Um GIF gravado (usando ferramentas como *ScreenToGif* ou *Loom*) mostrando o usuário digitando uma pergunta sobre estoque ou pedido e a IA respondendo em tempo real.

### 3. Modos Claro e Escuro (Lado a Lado)
- **Onde:** Na seção de funcionalidades.
- **O que colocar:** Uma tabela comparando o visual no tema *Light* e *Dark* da interface.

```markdown
| Modo Claro | Modo Escuro |
| :---: | :---: |
| <img src="docs/screen-light.png" width="400"> | <img src="docs/screen-dark.png" width="400"> |
```

### 4. Diagrama de Arquitetura

- **Onde:** Na seção de Métodos e Arquitetura.
- **O que colocar:** Um diagrama criado no *Mermaid.js* ou *Excalidraw* mostrando a relação entre o **Front-end Next.js 🚀** → **Back-end Go/GORM 🦫** → **PostgreSQL 🐘** → **API Gemini 🤖**.

## 🚀 Como Executar o Projeto

1. **Clone o repositório:**

   Bash
   ```
   git clone [https://github.com/seu-usuario/winthor-agent-ai.git](https://github.com/seu-usuario/winthor-agent-ai.git)
   cd winthor-agent-ai
   ```

2. **Instale as dependências:**

   Bash
   ```
   npm install
   ```

3. **Configure as Variáveis de Ambiente (****`.env.local`****):**

   Snippet de código
   ```
   NEXT_PUBLIC_API_URL=http://localhost:8080
   JWT_SECRET=sua_chave_secreta
   ```

4. **Inicie o servidor de desenvolvimento:**

   Bash
   ```
   npm run dev
   ```
