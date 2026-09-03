import { ChatShell } from "@/components/layout/chat/chat-shell";
import { getTokenPayload } from "@/service/TokenValid";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Chat - IA COMAL",
  description: "Chat with our AI assistant powered by Gemini",
  icons: {
    icon: [
      {
        url: "/winthor-ia/favicon.svg", // Ícone exibido no Modo Claro
        media: "(prefers-color-scheme: light)",
      },
      {
        url: "/winthor-ia/favicon-black.svg", // Ícone exibido no Modo Escuro
        media: "(prefers-color-scheme: dark)",
      },
    ],
  },
};

export default async function ChatPage() {
  // 1. Recupera os dados do token no Server Component
  const token = await getTokenPayload();

  // 2. Extrai o nome do usuário (com fallback para login/idPg caso necessário)
  const currentUserId =
    token?.nome || token?.usuario || String(token?.idPg || "Usuário");

  // 3. Repassa o userId para o ChatShell
  return <ChatShell userId={currentUserId} />;
}
