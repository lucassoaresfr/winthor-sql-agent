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
  const token = await getTokenPayload();
  const currentUserId =
    token?.nome || token?.usuario || String(token?.idPg || "Usuário");
    
  return <ChatShell userId={currentUserId} />;
}
