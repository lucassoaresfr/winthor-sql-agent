"use client";

import { useEffect, useState, useTransition, useCallback } from "react";
import { useRouter, useParams } from "next/navigation";
import { Plus, MessageSquare, Loader2 } from "lucide-react";

import { chatService } from "@/service/sidebar/route";
import { ChatResponse } from "@/types/sidebar";
import { Button } from "@/components/ui/button";
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarFooter,
} from "@/components/ui/sidebar";
import Image from "next/image";
import imagelogo from "@/../public/favicon-black.svg";

interface AppSidebarProps extends React.ComponentProps<typeof Sidebar> {
  userId: string; // Recebe o NOME do usuário
  userLogin?: string; // Login/matrícula opcional para o subtexto
}

export function AppSidebar({
  userId,
  userLogin,
  variant = "sidebar",
  ...props
}: AppSidebarProps) {
  const router = useRouter();
  const params = useParams();

  const [chats, setChats] = useState<ChatResponse[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [isPending, startTransition] = useTransition();

  const currentChatId = params?.id as string;

  // Extrai as iniciais a partir do NOME enviado em userId
  const getInitials = (name: string) => {
    if (!name) return "US";
    const parts = name.trim().split(" ");
    if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase();
    return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
  };

  const loadChats = useCallback(async () => {
    if (!userId) return;
    try {
      setLoading(true);
      const data = await chatService.listUserChats(userId);
      setChats(data || []);
    } catch (error) {
      console.error("Erro ao carregar lista de chats:", error);
    } finally {
      setLoading(false);
    }
  }, [userId]);

  useEffect(() => {
    loadChats();

    const handleChatCreated = () => {
      loadChats();
    };

    window.addEventListener("chat-created", handleChatCreated);

    return () => {
      window.removeEventListener("chat-created", handleChatCreated);
    };
  }, [userId, currentChatId, loadChats]);

  const handleCreateNewChat = () => {
    startTransition(() => {
      router.push("/");
    });
  };

  return (
    <Sidebar variant={variant} collapsible="icon" {...props}>
      <SidebarHeader className="p-3 group-data-[collapsible=icon]:p-2 border-b border-sidebar-border">
        {/* Logo / Título */}
        <div className="flex items-center gap-2 group-data-[collapsible=icon]:justify-center">
          <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-emerald-600 text-primary-foreground font-bold">
            <Image
              src={imagelogo}
              alt="Logo Winthor Agent"
              width={20}
              height={20}
              className="object-contain"
            />
          </div>
          <div className="flex flex-col gap-0.5 leading-none group-data-[collapsible=icon]:hidden">
            <span className="font-semibold text-sm">Winthor Agent</span>
            <span className="text-xs text-muted-foreground">Assistente IA</span>
          </div>
        </div>

        {/* Botão Novo Chat */}
        <div className="mt-2 flex justify-center">
          <Button
            onClick={handleCreateNewChat}
            disabled={isPending}
            className="w-full group-data-[collapsible=icon]:w-8 group-data-[collapsible=icon]:h-8 group-data-[collapsible=icon]:p-0 justify-center gap-2 shadow-sm bg-emerald-600 hover:bg-emerald-800"
            title="Nova Conversa"
          >
            {isPending ? (
              <Loader2 className="size-4 animate-spin shrink-0" />
            ) : (
              <Plus className="size-4 shrink-0" />
            )}
            <span className="group-data-[collapsible=icon]:hidden">
              Nova Conversa
            </span>
          </Button>
        </div>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel className="group-data-[collapsible=icon]:hidden">
            Histórico de Conversas
          </SidebarGroupLabel>
          <SidebarGroupContent>
            {loading ? (
              <div className="flex items-center justify-center py-6 text-muted-foreground">
                <Loader2 className="size-5 animate-spin" />
              </div>
            ) : chats.length === 0 ? (
              <div className="p-4 text-center text-xs text-muted-foreground group-data-[collapsible=icon]:hidden">
                Nenhum histórico encontrado.
              </div>
            ) : (
              <SidebarMenu>
                {chats.map((chat) => {
                  const isActive = currentChatId === chat.id;

                  return (
                    <SidebarMenuItem key={chat.id}>
                      <SidebarMenuButton
                        isActive={isActive}
                        onClick={() => router.push(`/chat/${chat.id}`)}
                        tooltip={chat.title || "Conversa sem título"}
                      >
                        <MessageSquare className="size-4 shrink-0" />
                        <span>{chat.title || "Nova Conversa"}</span>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  );
                })}
              </SidebarMenu>
            )}
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      {/* Rodapé exibindo Nome (userId) + Avatar com as Iniciais */}
      <SidebarFooter className="border-t border-sidebar-border p-3 group-data-[collapsible=icon]:p-2">
        <div className="flex items-center gap-3 group-data-[collapsible=icon]:justify-center">
          {/* Avatar com as Iniciais do Nome */}
          <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-emerald-100 text-emerald-800 font-semibold text-xs border border-emerald-300">
            {getInitials(userId)}
          </div>

          {/* Exibição do Nome Completo */}
          <div className="flex flex-col min-w-0 leading-tight group-data-[collapsible=icon]:hidden">
            <span
              className="font-medium text-sm truncate text-stone-800"
              title={userId}
            >
              {userId}
            </span>
            {userLogin && (
              <span className="text-[11px] text-muted-foreground truncate">
                @{userLogin}
              </span>
            )}
          </div>
        </div>
      </SidebarFooter>
    </Sidebar>
  );
}
