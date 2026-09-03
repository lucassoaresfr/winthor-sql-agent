// ==========================================
// TIPOS BASE / ENTIDADES
// ==========================================

export type MessageRole = "user" | "assistant" | "system";

export interface MessageResponse {
  id: string;
  chat_id: string;
  role: MessageRole;
  content: string;
  model_used?: string;
  created_at: string; // ISO 8601 string vinda da API
}

export interface ChatResponse {
  id: string;
  user_id: string;
  title: string;
  created_at: string;
  updated_at: string;
  messages?: MessageResponse[];
}

// ==========================================
// DTOS POR ROTA (REQUESTS & RESPONSES)
// ==========================================

/**
 * 1. POST /api/v1/chats
 * Criar uma nova sessão de conversa
 */
export interface CreateChatRequest {
  user_id: string;
  title?: string;
}

export type CreateChatResponse = ChatResponse;

/**
 * 2. POST /api/v1/chats/:id/messages
 * Adicionar uma mensagem ao histórico
 */
export interface SaveMessageRequest {
  role: MessageRole;
  content: string;
  model_used?: string;
  prompt_tokens?: number;
  completion_tokens?: number;
}

export type SaveMessageResponse = MessageResponse;

/**
 * 3. GET /api/v1/chats/:id
 * Obter histórico completo de um chat pelo ID
 */
export interface GetChatByIdParams {
  id: string;
}

export type GetChatByIdResponse = ChatResponse;

/**
 * 4. GET /api/v1/chats/user/:user_id
 * Listar todas as conversas de um usuário (Sidebar)
 */
export interface ListUserChatsParams {
  user_id: string;
}

export type ListUserChatsResponse = ChatResponse[];
