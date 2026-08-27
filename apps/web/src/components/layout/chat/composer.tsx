"use client";

import type React from "react";
import {
  useState,
  useRef,
  useCallback,
  type KeyboardEvent,
  useEffect,
} from "react";
import Image from "next/image";
import { cn } from "@/lib/utils";
import { AnimatedOrb } from "./animated-orb";
import { Brain } from "lucide-react";

interface ComposerProps {
  onSend: (content: string) => void;
  isStreaming: boolean;
  disabled?: boolean;
}

export function Composer({ onSend, isStreaming, disabled }: ComposerProps) {
  const [value, setValue] = useState("");
  const [hasAnimated, setHasAnimated] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    setHasAnimated(true);
  }, []);

  const playClickSound = useCallback(() => {
    const audio = new Audio(
      "https://hebbkx1anhila5yf.public.blob.vercel-storage.com/click-FM4Xaa1FJj237591TiZw4yL1fIxdOw.mp3",
    );
    audio.volume = 0.5;
    audio.play().catch(() => {});
  }, []);

  const handleInput = useCallback(() => {
    const textarea = textareaRef.current;
    if (textarea) {
      textarea.style.height = "auto";
      textarea.style.height = `${Math.min(textarea.scrollHeight, 200)}px`;
    }
  }, []);

  const handleSend = useCallback(() => {
    if (!value.trim() || isStreaming || disabled) return;
    playClickSound();

    onSend(value.trim());
    setValue("");
    if (textareaRef.current) {
      textareaRef.current.style.height = "auto";
    }
  }, [value, isStreaming, disabled, onSend, playClickSound]);

  const handleKeyDown = useCallback(
    (e: KeyboardEvent<HTMLTextAreaElement>) => {
      if (e.key === "Enter" && !e.shiftKey) {
        e.preventDefault();
        handleSend();
      }
    },
    [handleSend],
  );

  return (
    <div
      className={cn(
        "fixed bottom-4 left-0 right-0 px-4 pointer-events-none z-10",
        hasAnimated && "composer-intro",
      )}
    >
      <div className="relative max-w-2xl mx-auto pointer-events-auto">
        <div
          className={cn(
            "flex flex-col gap-3 p-4 bg-white border-stone-200 transition-all duration-200 border-none overflow-hidden relative rounded-3xl",
            "focus-within:border-stone-300 focus-within:ring-2 focus-within:ring-stone-200",
          )}
          style={{
            boxShadow:
              "rgba(14, 63, 126, 0.06) 0px 0px 0px 1px, rgba(42, 51, 69, 0.06) 0px 1px 1px -0.5px, rgba(42, 51, 70, 0.06) 0px 3px 3px -1.5px, rgba(42, 51, 70, 0.06) 0px 6px 6px -3px, rgba(14, 63, 126, 0.06) 0px 12px 12px -6px, rgba(14, 63, 126, 0.06) 0px 24px 24px -12px",
          }}
        >
          <div className="flex gap-2 items-center">
            <textarea
              ref={textareaRef}
              value={value}
              onChange={(e) => {
                setValue(e.target.value);
                handleInput();
              }}
              onKeyDown={handleKeyDown}
              placeholder="Digite uma mensagem... (Shift+Enter para nova linha)"
              disabled={isStreaming || disabled}
              rows={1}
              className={cn(
                "flex-1 resize-none bg-transparent px-2 py-1.5 text-sm text-stone-800 placeholder:text-stone-400",
                "focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed",
                "max-h-14 overflow-y-auto",
              )}
              aria-label="Entrada de mensagem"
            />

            <button
              onClick={handleSend}
              disabled={!value.trim() || isStreaming || disabled}
              className={cn(
                "relative h-9 w-9 shrink-0 transition-all rounded-full flex items-center justify-center",
                !value.trim() || isStreaming || disabled
                  ? "opacity-50 cursor-not-allowed"
                  : "cursor-pointer hover:scale-105",
              )}
              aria-label="Enviar mensagem"
            >
              <AnimatedOrb size={36} />
            </button>
          </div>

          <div className="flex items-center gap-2 px-1 pt-1">
            <div className="flex items-center gap-1.5 px-2 py-1 rounded-full bg-zinc-100 text-stone-600">
              <Brain size={16}/>
              <span className="text-xs font-medium">Gemini (auto)</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
