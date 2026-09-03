"use client";

import { cn } from "@/lib/utils";
import type { Message } from "./chat-shell";
import { User } from "lucide-react";
import { MarkdownRenderer } from "./markdown-renderer";
import Image from "next/image";
import { AnimatedOrb } from "./animated-orb";

interface MessageBubbleProps {
  message: Message;
  isStreaming?: boolean;
}

function formatTime(date: Date): string {
  return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
}

export function MessageBubble({
  message,
  isStreaming = false,
}: MessageBubbleProps) {
  const isUser = message.role === "user";
  const hasContent = message.content && message.content.trim().length > 0;

  return (
    <div
      className={cn(
        "flex max-w-[90%] md:max-w-[80%] gap-3.5",
        isUser
          ? "ml-auto flex-row-reverse user-message-enter"
          : "mr-auto animate-in fade-in slide-in-from-bottom-2 duration-300 items-start",
      )}
    >
      {/* Avatar */}
      <div
        className={cn(
          "w-8 h-8 rounded-full flex items-center justify-center shrink-0 mt-1",
          isUser ? "bg-stone-100 border border-stone-200" : "bg-transparent",
        )}
        style={{
          boxShadow: isUser ? "0px 2px 4px rgba(0,0,0,0.04)" : "none",
        }}
        aria-hidden="true"
      >
        {isUser ? (
          <User className="w-4 h-4 text-stone-700" />
        ) : (
          <AnimatedOrb className="w-8 h-8 shrink-0" />
        )}
      </div>

      {/* Conteúdo da mensagem */}
      <div
        className={cn("flex flex-col", isUser ? "items-end" : "items-start")}
      >
        <span className="text-[11px] font-medium text-stone-400 mb-1 px-1">
          {isUser ? "Você" : "Assistente"}
        </span>

        <div
          className={cn(
            "rounded-2xl overflow-hidden transition-all duration-300",
            isUser
              ? "bg-emerald-600 text-white rounded-br-xs shadow-sm"
              : "bg-white/90 backdrop-blur-sm text-stone-800 border border-stone-100 rounded-bl-xs shadow-sm",
          )}
        >
          <div
            className={cn(
              "px-4 py-3 relative transition-all duration-300",
              !isUser && isStreaming && "pr-6",
            )}
          >
            {isUser ? (
              <div className="flex flex-col gap-2">
                {message.imageData && (
                  <div className="w-20 h-20 rounded-lg overflow-hidden border border-emerald-500/30">
                    <Image
                      src={message.imageData || "/placeholder.svg"}
                      alt="Uploaded image"
                      width={80}
                      height={80}
                      className="w-full h-full object-cover"
                    />
                  </div>
                )}
                <p className="text-sm whitespace-pre-wrap wrap-break-words">
                  {message.content}
                </p>
              </div>
            ) : (
              <div className="relative text-sm">
                {!hasContent ? (
                  /* Transição: 3 pontos alinhados que aparecem assim que o assistente assume */
                  <div className="flex items-center gap-1.5 py-1 px-0.5">
                    <span className="w-2 h-2 bg-emerald-500 rounded-full animate-ping animation-duration-[1s]" />
                    <span className="w-2 h-2 bg-emerald-500 rounded-full animate-ping animation-duration-[1s] [animation-delay:200ms]" />
                    <span className="w-2 h-2 bg-emerald-500 rounded-full animate-ping animation-duration-[1s] [animation-delay:400ms]" />
                  </div>
                ) : (
                  <>
                    <MarkdownRenderer content={message.content} />
                    {/* Cursor pulsante no final do streaming */}
                    {isStreaming && (
                      <span className="inline-block w-2 h-4 ml-1 bg-emerald-500 rounded-xs animate-pulse align-middle" />
                    )}
                  </>
                )}
              </div>
            )}
          </div>
        </div>

        <span className="text-[10px] text-stone-400 mt-1 px-1">
          {formatTime(message.createdAt)}
        </span>
      </div>
    </div>
  );
}
