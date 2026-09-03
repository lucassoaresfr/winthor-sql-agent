"use client";

import { useState, useEffect, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import { MessageSquareDashed, Loader2 } from "lucide-react";
import { MessageList } from "./message-list";
import { Composer } from "./composer";
import { Button } from "@/components/ui/button";
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
  const chatId = params?.id as string | undefined;

  const [messages, setMessages] = useState<Message[]>([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isLoaded, setIsLoaded] = useState(false);

  // 1. Carregar histórico do chat do PostgreSQL ao acessar a rota /chat/[id]
  useEffect(() => {
    async function loadChatHistory() {
      if (!chatId) {
        setMessages([]);
        setIsLoaded(true);
        setError(null);
        return;
      }

      try {
        setIsLoaded(false);
        setError(null);
        const data = await chatService.getChatById(chatId);
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
  }, [chatId]);

  // 2. Enviar mensagem e persistir no PostgreSQL
  const sendMessage = useCallback(
    async (content: string) => {
      const trimmedContent = content.trim();
      if (!trimmedContent || isStreaming) return;

      setError(null);
      setIsStreaming(true);

      let activeChatId = chatId;

      try {
        // A. Se não existir um chat ativo na URL, cria a sessão no Postgres
        if (!activeChatId) {
          const newChat = await chatService.createChat({
            user_id: CURRENT_USER_ID,
            title:
              trimmedContent.slice(0, 30) +
              (trimmedContent.length > 30 ? "..." : ""),
          });
          activeChatId = newChat.id;

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

        const currentMessages = [...messages, userMessage];
        setMessages(currentMessages);

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
    [chatId, messages, isStreaming],
  );

  // 3. Função de Re-tentativa (Retry)
  const handleRetry = useCallback(async () => {
    if (isStreaming) return;

    // Se já existem mensagens na conversa, recupera a última enviada pelo usuário
    const lastUserMessage = [...messages]
      .reverse()
      .find((m) => m.role === "user");

    if (!lastUserMessage || !chatId) {
      setError(null);
      return;
    }

    setError(null);
    setIsStreaming(true);

    try {
      // Reenvia apenas a janela de contexto para a IA (sem re-salvar a msg do usuário no banco)
      const contextWindow = messages.slice(-MAX_CONTEXT_MESSAGES).map((m) => ({
        role: m.role,
        content: m.content,
      }));

      const response = await ChatApi({ messages: contextWindow });
      if (response.error) throw new Error(response.error);

      const respostaTexto =
        response.resposta || "Não foi possível obter uma resposta.";

      // Salva a resposta gerada da IA
      const assistantMsgSaved = await chatService.addMessage(chatId, {
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
  }, [chatId, messages, isStreaming]);

  // const clearChat = useCallback(() => {
  //   setMessages([]);
  //   setError(null);
  //   router.push("/");
  // }, [router]);

  return (
    <div className="relative flex flex-col h-full w-full bg-stone-50 overflow-hidden">
      {/* <Button
        onClick={clearChat}
        variant="ghost"
        size="icon"
        className="absolute top-4 left-4 z-20 h-10 w-10 rounded-full bg-zinc-100 hover:bg-zinc-200 text-stone-600 cursor-pointer shadow-sm"
        aria-label="Nova conversa"
        title="Nova conversa"
      >
        <MessageSquareDashed className="w-5 h-5" />
      </Button> */}

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
