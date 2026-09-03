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
} from "@/components/ui/sidebar";
import Image from "next/image";
import imagelogo from "@/../public/favicon-black.svg";

interface AppSidebarProps extends React.ComponentProps<typeof Sidebar> {
  userId: string;
}

export function AppSidebar({
  userId,
  variant = "sidebar",
  ...props
}: AppSidebarProps) {
  const router = useRouter();
  const params = useParams();

  const [chats, setChats] = useState<ChatResponse[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [isPending, startTransition] = useTransition();

  const currentChatId = params?.id as string;

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

    // Escuta o evento customizado disparado no ChatShell ao criar um novo chat
    const handleChatCreated = () => {
      loadChats();
    };

    window.addEventListener("chat-created", handleChatCreated);

    return () => {
      window.removeEventListener("chat-created", handleChatCreated);
    };
  }, [userId, currentChatId, loadChats]);

  // Navega para a raiz (/), deixando a criação da sessão para o envio da primeira mensagem
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
    </Sidebar>
  );
}
