"use client";

import { useEffect, useRef, useState } from "react";
import { MessageBubble } from "./message-bubble";
import type { Message } from "./chat-shell";
import { TypingIndicator } from "./typing-indicator";
import { AlertCircle, RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { AnimatedOrb } from "./animated-orb";

interface MessageListProps {
  messages: Message[];
  isStreaming: boolean;
  error: string | null;
  onRetry: () => void;
  isLoaded: boolean;
}

const LAUNCH_SOUND_URL =
  "https://hebbkx1anhila5yf.public.blob.vercel-storage.com/launch-SUi0itAGHr1wtvdDYYG5bzFLsIYHtP.mp3";

export function MessageList({
  messages,
  isStreaming,
  error,
  onRetry,
  isLoaded,
}: MessageListProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);
  const [autoScroll, setAutoScroll] = useState(true);
  const [hasAnimated, setHasAnimated] = useState(false);
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const hasPlayedIntroRef = useRef(false);

  useEffect(() => {
    if (!isLoaded) return;

    if (messages.length === 0 && !hasPlayedIntroRef.current) {
      setHasAnimated(true);
      hasPlayedIntroRef.current = true;

      audioRef.current = new Audio(LAUNCH_SOUND_URL);
      audioRef.current.volume = 0.5;
      audioRef.current.play().catch(() => {});
    } else if (messages.length > 0) {
      setHasAnimated(false);
      hasPlayedIntroRef.current = true;
    }

    return () => {
      if (audioRef.current) {
        audioRef.current.pause();
        audioRef.current = null;
      }
    };
  }, [isLoaded, messages.length]);

  useEffect(() => {
    if (!containerRef.current || !autoScroll) return;
    containerRef.current.scrollTo({
      top: containerRef.current.scrollHeight,
      behavior: isStreaming ? "smooth" : "auto",
    });
  }, [messages, isStreaming, autoScroll]);

  const handleScroll = () => {
    if (!containerRef.current) return;
    const { scrollTop, scrollHeight, clientHeight } = containerRef.current;
    const isAtBottom = scrollHeight - scrollTop - clientHeight < 100;
    setAutoScroll(isAtBottom);
  };

  const filteredMessages = messages.filter((message, index) => {
    const isLast = index === messages.length - 1;

    if (
      isStreaming &&
      message.role === "assistant" &&
      isLast &&
      !message.content
    ) {
      return false;
    }

    if (error && message.role === "assistant" && isLast) {
      return false;
    }

    return true;
  });

  const showTypingIndicator =
    isStreaming &&
    (filteredMessages.length === 0 ||
      filteredMessages[filteredMessages.length - 1]?.role === "user");

  if (!isLoaded) {
    return (
      <div className="flex flex-1 items-center justify-center h-full">
        <AnimatedOrb size={64} />
      </div>
    );
  }

  // ESTADO VAZIO: Centralizado sem rolagem
  if (filteredMessages.length === 0 && !error && !isStreaming) {
    return (
      <div className="flex flex-col items-center justify-center h-full text-center text-stone-400 p-4">
        <div
          className={`mb-4 transition-transform duration-500 ${
            hasAnimated ? "scale-105" : ""
          }`}
        >
          <AnimatedOrb size={110} />
        </div>
        <p className="text-lg font-semibold text-stone-800 tracking-tight">
          Olá! Sou a assistente inteligente da COMAL
        </p>
        <p className="text-sm mt-1.5 text-stone-500 max-w-sm leading-relaxed">
          Como posso ajudar hoje? Consulte pedidos, clientes, estoque ou
          relatórios da empresa.
        </p>
      </div>
    );
  }

  // ESTADO COM MENSAGENS
  return (
    <div
      ref={containerRef}
      onScroll={handleScroll}
      className="h-full overflow-y-auto py-6 space-y-4 px-4 md:px-8 scroll-smooth"
      role="log"
      aria-label="Mensagens do chat"
      aria-live="polite"
    >
      {filteredMessages.map((message) => (
        <MessageBubble
          key={message.id}
          message={message}
          isStreaming={
            isStreaming &&
            message.role === "assistant" &&
            message === filteredMessages[filteredMessages.length - 1]
          }
        />
      ))}

      {showTypingIndicator && <TypingIndicator />}

      {error && (
        <div className="flex items-center gap-3 p-4 bg-amber-50/80 border border-amber-200/60 rounded-xl shadow-xs">
          <AlertCircle className="w-5 h-5 text-amber-600 shrink-0" />
          <div className="flex-1">
            <p className="text-xs font-semibold text-amber-900">
              Aviso do Sistema
            </p>
            <p className="text-xs text-amber-700 mt-0.5">{error}</p>
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={onRetry}
            className="text-amber-800 hover:bg-amber-100/60 text-xs cursor-pointer"
          >
            <RefreshCw className="w-3.5 h-3.5 mr-1.5" />
            Tentar novamente
          </Button>
        </div>
      )}

      <div ref={bottomRef} className="h-4" />
    </div>
  );
}
