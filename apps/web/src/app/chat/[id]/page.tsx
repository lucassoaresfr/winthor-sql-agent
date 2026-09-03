import { ChatShell } from "@/components/layout/chat/chat-shell";
import { Metadata } from "next";

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
export default function ChatPage() {
  return <ChatShell />;
}
