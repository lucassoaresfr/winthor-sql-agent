export interface ChatApiMessage {
  role: "user" | "assistant";
  content: string;
}

export interface ChatApiRequest {
  messages: ChatApiMessage[];
}

export interface ChatApiResponse {
  resposta?: string;
  error?: string;
  status?: string;
}
