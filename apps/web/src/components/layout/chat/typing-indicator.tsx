"use client";

import { AnimatedOrb } from "./animated-orb";

export function TypingIndicator() {
  return (
    <div className="flex items-center gap-3 mr-auto animate-in fade-in slide-in-from-bottom-2 duration-300">
      <div className="w-8 h-8 rounded-full flex items-center justify-center shrink-0">
        <AnimatedOrb className="w-8 h-8" />
      </div>
      <div className="bg-white/90 backdrop-blur-md border border-stone-200/60 rounded-2xl rounded-bl-xs px-4 py-3 shadow-xs flex items-center justify-center">
        {/* Contêiner em rotação circular contínua */}
        <div className="relative w-5 h-5 animate-spin-orbit">
          <span className="absolute top-0 left-1/2 -translate-x-1/2 w-1.5 h-1.5 bg-emerald-500 rounded-full animate-pulse" />
          <span className="absolute bottom-0 left-0 w-1.5 h-1.5 bg-emerald-500 rounded-full animate-pulse [animation-delay:200ms]" />
          <span className="absolute bottom-0 right-0 w-1.5 h-1.5 bg-emerald-500 rounded-full animate-pulse [animation-delay:400ms]" />
        </div>
      </div>

      <style jsx>{`
        @keyframes spinOrbit {
          0% {
            transform: rotate(0deg);
          }
          100% {
            transform: rotate(360deg);
          }
        }
        .animate-spin-orbit {
          animation: spinOrbit 1.1s cubic-bezier(0.5, 0.1, 0.5, 0.9) infinite;
        }
      `}</style>
    </div>
  );
}
