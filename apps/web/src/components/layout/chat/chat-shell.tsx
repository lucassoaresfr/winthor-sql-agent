"use client";

import { useState, useEffect, useCallback } from "react";
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
const WELCOME_CONTENT = `👋 **Olá! Sou o assistente virtual do WinThor.**

Posso te ajudar a consultar dados e informações do sistema em tempo real:

📦 **Produtos e Estoque:** Preços, saldos para venda e categorias (*Departamento/Seção*).

🛒 **Pedidos de Venda:** Status (*Faturado, Liberado, Pendente*), valores e itens do pedido.

🏷️ **Promoções:** Descontos do dia, ofertas ativas e itens em promoção.

🏢 **Clientes e CNPJ:** Consultas cadastrais internas e busca pública na Receita Federal.

💬 *Como posso te ajudar hoje?*`;

const VIRTUAL_WELCOME_MESSAGE: Message = {
  id: "welcome-initial",
  role: "assistant",
  content: WELCOME_CONTENT,
  createdAt: new Date(),
};

export function ChatShell() {
  const params = useParams();
  const chatIdParam = params?.id as string | undefined;

  const [activeChatId, setActiveChatId] = useState<string | undefined>(
    chatIdParam,
  );

  const [messages, setMessages] = useState<Message[]>([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isLoaded, setIsLoaded] = useState(false);

  useEffect(() => {
    setActiveChatId(chatIdParam);
  }, [chatIdParam]);

  // 1. Carregar histórico do chat do PostgreSQL ao acessar /chat/[id]
  useEffect(() => {
    async function loadChatHistory() {
      if (!activeChatId) {
        setMessages([]);
        setIsLoaded(true);
        setError(null);
        return;
      }

      try {
        setIsLoaded(false);
        setError(null);
        const data = await chatService.getChatById(activeChatId);
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
  }, [activeChatId]);

  // 2. Enviar mensagem e persistir no PostgreSQL
  const sendMessage = useCallback(
    async (content: string) => {
      const trimmedContent = content.trim();
      if (!trimmedContent || isStreaming) return;

      setError(null);
      setIsStreaming(true);

      let currentChatId = activeChatId;
      let initialMessagesList: Message[] = [...messages];

      try {
        // A. Se não existir um chat ativo, cria a sessão e grava a mensagem de boas-vindas no banco
        if (!currentChatId) {
          const newChat = await chatService.createChat({
            user_id: CURRENT_USER_ID,
            title:
              trimmedContent.slice(0, 30) +
              (trimmedContent.length > 30 ? "..." : ""),
          });
          currentChatId = newChat.id;
          setActiveChatId(currentChatId);

          // 🟢 1. Grava a mensagem de boas-vindas do assistente no banco
          const welcomeSaved = await chatService.addMessage(currentChatId, {
            role: "assistant",
            content: WELCOME_CONTENT,
          });

          const welcomeMsgObj: Message = {
            id: welcomeSaved.id,
            role: "assistant",
            content: welcomeSaved.content,
            createdAt: new Date(welcomeSaved.created_at),
          };

          initialMessagesList = [welcomeMsgObj];

          // Atualiza a URL e a sidebar
          window.history.replaceState(null, "", `/chat/${currentChatId}`);
          window.dispatchEvent(new Event("chat-created"));
        }

        // B. Salvar a mensagem do usuário no banco de dados
        const userMsgSaved = await chatService.addMessage(currentChatId, {
          role: "user",
          content: trimmedContent,
        });

        const userMessage: Message = {
          id: userMsgSaved.id,
          role: "user",
          content: userMsgSaved.content,
          createdAt: new Date(userMsgSaved.created_at),
        };

        const updatedMessages = [...initialMessagesList, userMessage];
        setMessages(updatedMessages);

        // C. Processar resposta da IA com o contexto COMPLETO (incluindo a boas-vindas)
        const contextWindow = updatedMessages
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
        const assistantMsgSaved = await chatService.addMessage(currentChatId, {
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
    [activeChatId, messages, isStreaming],
  );

  // 3. Função de Re-tentativa (Retry)
  const handleRetry = useCallback(async () => {
    if (isStreaming) return;

    const lastUserMessage = [...messages]
      .reverse()
      .find((m) => m.role === "user");

    if (!lastUserMessage || !activeChatId) {
      setError(null);
      return;
    }

    setError(null);
    setIsStreaming(true);

    try {
      const contextWindow = messages.slice(-MAX_CONTEXT_MESSAGES).map((m) => ({
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
  }, [activeChatId, messages, isStreaming]);

  // 🟢 Se estiver na rota inicial (sem chat ativo) e sem mensagens, exibe a mensagem de boas-vindas virtualmente
  const displayMessages =
    messages.length === 0 && !activeChatId
      ? [VIRTUAL_WELCOME_MESSAGE]
      : messages;

  return (
    <div className="relative flex flex-col h-full w-full bg-stone-50 overflow-hidden">
      {/* Área de mensagens */}
      <div
        className={`flex-1 ${
          displayMessages.length > 0 ? "overflow-y-auto" : "overflow-hidden"
        }`}
      >
        {!isLoaded ? (
          <div className="flex h-full items-center justify-center text-muted-foreground">
            <Loader2 className="size-6 animate-spin" />
          </div>
        ) : (
          <MessageList
            messages={displayMessages}
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
