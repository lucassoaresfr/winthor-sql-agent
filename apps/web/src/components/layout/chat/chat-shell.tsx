"use client";

import { useState, useEffect, useCallback } from "react";
import { MessageSquareDashed } from "lucide-react";
import { MessageList } from "./message-list";
import { Composer } from "./composer";
import { Button } from "@/components/ui/button";
import { ChatApi } from "@/service/chat/routes";

export interface Message {
  id: string;
  role: "user" | "assistant";
  content: string;
  createdAt: Date;
}

const STORAGE_KEY = "chat-messages";
// Quantidade ideal de mensagens recentes para manter o contexto sem estourar o limite de tokens
const MAX_CONTEXT_MESSAGES = 6;

function generateId(): string {
  return `${Date.now()}-${Math.random().toString(36).substring(2, 9)}`;
}

export function ChatShell() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isLoaded, setIsLoaded] = useState(false);

  // Carrega histórico do localStorage
  useEffect(() => {
    try {
      const stored = localStorage.getItem(STORAGE_KEY);
      if (stored) {
        const parsed = JSON.parse(stored);
        const messagesWithDates = parsed.map((msg: Message) => ({
          ...msg,
          createdAt: new Date(msg.createdAt),
        }));
        setMessages(messagesWithDates);
      }
    } catch (e) {
      console.error("Failed to load from localStorage:", e);
    } finally {
      setIsLoaded(true);
    }
  }, []);

  // Salva no localStorage
  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(messages));
    } catch (e) {
      console.error("Failed to save messages to localStorage:", e);
    }
  }, [messages]);

  // Envio de mensagens via ChatApi
  const sendMessage = useCallback(
    async (content: string) => {
      const trimmedContent = content.trim();
      if (!trimmedContent || isStreaming) return;

      setError(null);

      const userMessage: Message = {
        id: generateId(),
        role: "user",
        content: trimmedContent,
        createdAt: new Date(),
      };

      // Cria a lista atualizada de mensagens incluindo a nova pergunta
      const updatedMessages = [...messages, userMessage];
      setMessages(updatedMessages);
      setIsStreaming(true);

      // Prepara o payload recortando apenas as últimas N mensagens (Janela Deslizante)
      const contextWindow = updatedMessages
        .slice(-MAX_CONTEXT_MESSAGES)
        .map((m) => ({
          role: m.role,
          content: m.content,
        }));

      try {
        // Envia o array de mensagens para a API
        const response = await ChatApi({ messages: contextWindow });

        if (response.error) {
          throw new Error(response.error);
        }

        const respostaTexto =
          response.resposta || "Não foi possível obter uma resposta.";

        const assistantMessage: Message = {
          id: generateId(),
          role: "assistant",
          content: respostaTexto,
          createdAt: new Date(),
        };

        setMessages((prev) => [...prev, assistantMessage]);
      } catch (e: any) {
        console.error("Error sending message:", e);
        setError(e.message || "Ocorreu um erro ao processar sua pergunta.");
      } finally {
        setIsStreaming(false);
      }
    },
    [messages, isStreaming],
  );

  const retry = useCallback(() => {
    if (messages.length === 0) return;
    const lastUserMessage = [...messages]
      .reverse()
      .find((m) => m.role === "user");

    if (lastUserMessage) {
      const index = messages.findIndex((m) => m.id === lastUserMessage.id);
      setMessages(messages.slice(0, index));
      setError(null);
      setTimeout(() => sendMessage(lastUserMessage.content), 100);
    }
  }, [messages, sendMessage]);

  const clearChat = useCallback(() => {
    setMessages([]);
    setError(null);
    localStorage.removeItem(STORAGE_KEY);
  }, []);

  return (
    <div
      className="relative h-dvh bg-stone-50"
      style={{
        boxShadow:
          "rgba(14, 63, 126, 0.04) 0px 0px 0px 1px, rgba(42, 51, 69, 0.04) 0px 1px 1px -0.5px, rgba(42, 51, 70, 0.04) 0px 3px 3px -1.5px, rgba(42, 51, 70, 0.04) 0px 6px 6px -3px, rgba(14, 63, 126, 0.04) 0px 12px 12px -6px, rgba(14, 63, 126, 0.04) 0px 24px 24px -12px",
      }}
    >
      <Button
        onClick={clearChat}
        variant="ghost"
        size="icon"
        className="absolute top-4 left-4 z-20 h-10 w-10 rounded-full bg-zinc-100 hover:bg-zinc-200 text-stone-600"
        aria-label="Reset chat"
      >
        <MessageSquareDashed className="w-5 h-5" />
      </Button>

      <MessageList
        messages={messages}
        isStreaming={isStreaming}
        error={error}
        onRetry={retry}
        isLoaded={isLoaded}
      />

      <Composer
        onSend={sendMessage}
        isStreaming={isStreaming}
        disabled={!!error}
      />
    </div>
  );
}
