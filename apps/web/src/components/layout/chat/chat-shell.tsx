"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import { useParams } from "next/navigation";
import { Loader2 } from "lucide-react";
import { MessageList } from "./message-list";
import { Composer } from "./composer";
import { chatService } from "@/service/sidebar/route";
import { ChatApi } from "@/service/chat/routes";

export interface Message {
  id: string;
  role: "user" | "assistant";
  content: string;
  createdAt: Date;
  imageData?: string;
}

const CURRENT_USER_ID = "user";
const MAX_CONTEXT_MESSAGES = 6;

export function ChatShell() {
  const params = useParams();
  const urlChatId = params?.id as string | undefined;

  const [messages, setMessages] = useState<Message[]>([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isLoaded, setIsLoaded] = useState(false);

  // 🔴 REF para garantir que o ID do chat ativo seja lido sem depender de re-render
  const currentChatIdRef = useRef<string | undefined>(urlChatId);

  // Mantém a ref sincronizada sempre que a URL mudar
  useEffect(() => {
    currentChatIdRef.current = urlChatId;
  }, [urlChatId]);

  // 1. Carregar histórico do chat do PostgreSQL ao acessar a rota /chat/[id]
  useEffect(() => {
    async function loadChatHistory() {
      if (!urlChatId) {
        setMessages([]);
        setIsLoaded(true);
        setError(null);
        return;
      }

      try {
        setIsLoaded(false);
        setError(null);
        const data = await chatService.getChatById(urlChatId);
        if (data && data.messages) {
          const formattedMessages: Message[] = data.messages.map((m: any) => ({
            id: m.id,
            role: m.role as "user" | "assistant",
            content: m.content,
            createdAt: new Date(m.created_at),
          }));
          setMessages(formattedMessages);
        }
      } catch (err) {
        console.error("Erro ao carregar histórico do chat:", err);
        setError("Não foi possível carregar o histórico desta conversa.");
      } finally {
        setIsLoaded(true);
      }
    }

    loadChatHistory();
  }, [urlChatId]);

  // 2. Enviar mensagem e persistir no PostgreSQL
  const sendMessage = useCallback(
    async (content: string) => {
      const trimmedContent = content.trim();
      if (!trimmedContent || isStreaming) return;

      setError(null);
      setIsStreaming(true);

      // Usa a REF para pegar o chatId mais recente atualizado
      let activeChatId = currentChatIdRef.current;

      try {
        // A. Se não existir um chat ativo, cria a sessão no Postgres
        if (!activeChatId) {
          const newChat = await chatService.createChat({
            user_id: CURRENT_USER_ID,
            title:
              trimmedContent.slice(0, 30) +
              (trimmedContent.length > 30 ? "..." : ""),
          });

          activeChatId = newChat.id;
          currentChatIdRef.current = newChat.id; // 🟢 Atualiza a REF imediatamente!

          // 1. Atualiza a URL sem recarregar a página
          window.history.replaceState(null, "", `/chat/${activeChatId}`);

          // 2. Dispara evento para avisar a Sidebar para recarregar o histórico
          window.dispatchEvent(new Event("chat-created"));
        }

        // B. Salvar a mensagem do usuário no banco de dados
        const userMsgSaved = await chatService.addMessage(activeChatId, {
          role: "user",
          content: trimmedContent,
        });

        const userMessage: Message = {
          id: userMsgSaved.id,
          role: "user",
          content: userMsgSaved.content,
          createdAt: new Date(userMsgSaved.created_at),
        };

        // Usa atualização funcional para não depender de stale state em 'messages'
        let currentMessages: Message[] = [];
        setMessages((prev) => {
          currentMessages = [...prev, userMessage];
          return currentMessages;
        });

        // C. Processar resposta da IA
        const contextWindow = currentMessages
          .slice(-MAX_CONTEXT_MESSAGES)
          .map((m) => ({
            role: m.role,
            content: m.content,
          }));

        const response = await ChatApi({ messages: contextWindow });
        if (response.error) throw new Error(response.error);

        const respostaTexto =
          response.resposta || "Não foi possível obter uma resposta.";

        // D. Salvar a resposta da IA no banco de dados
        const assistantMsgSaved = await chatService.addMessage(activeChatId, {
          role: "assistant",
          content: respostaTexto,
        });

        const assistantMessage: Message = {
          id: assistantMsgSaved.id,
          role: "assistant",
          content: assistantMsgSaved.content,
          createdAt: new Date(assistantMsgSaved.created_at),
        };

        setMessages((prev) => [...prev, assistantMessage]);
      } catch (e: any) {
        console.error("Erro no envio da mensagem:", e);
        setError(e.message || "Ocorreu um erro ao processar sua pergunta.");
      } finally {
        setIsStreaming(false);
      }
    },
    [isStreaming], // Removidos chatId e messages das dependências para evitar closures antigas
  );

  // 3. Função de Re-tentativa (Retry)
  const handleRetry = useCallback(async () => {
    if (isStreaming) return;

    const activeChatId = currentChatIdRef.current;
    if (!activeChatId) {
      setError(null);
      return;
    }

    setError(null);
    setIsStreaming(true);

    try {
      let currentMessages: Message[] = [];
      setMessages((prev) => {
        currentMessages = prev;
        return prev;
      });

      const contextWindow = currentMessages
        .slice(-MAX_CONTEXT_MESSAGES)
        .map((m) => ({
          role: m.role,
          content: m.content,
        }));

      const response = await ChatApi({ messages: contextWindow });
      if (response.error) throw new Error(response.error);

      const respostaTexto =
        response.resposta || "Não foi possível obter uma resposta.";

      const assistantMsgSaved = await chatService.addMessage(activeChatId, {
        role: "assistant",
        content: respostaTexto,
      });

      const assistantMessage: Message = {
        id: assistantMsgSaved.id,
        role: "assistant",
        content: assistantMsgSaved.content,
        createdAt: new Date(assistantMsgSaved.created_at),
      };

      setMessages((prev) => [...prev, assistantMessage]);
    } catch (e: any) {
      console.error("Erro ao tentar novamente:", e);
      setError(e.message || "Ocorreu um erro ao processar sua pergunta.");
    } finally {
      setIsStreaming(false);
    }
  }, [isStreaming]);

  return (
    <div className="relative flex flex-col h-full w-full bg-stone-50 overflow-hidden">
      {/* Área de mensagens */}
      <div
        className={`flex-1 ${
          messages.length > 0 ? "overflow-y-auto" : "overflow-hidden"
        }`}
      >
        {!isLoaded ? (
          <div className="flex h-full items-center justify-center text-muted-foreground">
            <Loader2 className="size-6 animate-spin" />
          </div>
        ) : (
          <MessageList
            messages={messages}
            isStreaming={isStreaming}
            error={error}
            onRetry={handleRetry}
            isLoaded={isLoaded}
          />
        )}
      </div>

      {/* Rodapé do Composer */}
      <div className="p-3 sm:p-4 shrink-0 bg-stone-50/80 backdrop-blur-sm">
        <Composer onSend={sendMessage} isStreaming={isStreaming} />
      </div>
    </div>
  );
}
