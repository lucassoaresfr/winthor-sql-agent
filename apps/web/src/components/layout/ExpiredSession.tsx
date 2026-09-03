"use client";

import { useEffect } from "react";
import { Loader2, ShieldAlert } from "lucide-react"; // Use o ícone de sua preferência

export function ExpiredSession() {
  useEffect(() => {
    const timer = setTimeout(() => {
      window.location.href = "http://painelcomal:3000";
    }, 2500);

    return () => clearTimeout(timer);
  }, []);

  return (
    <div className="flex h-screen w-full items-center justify-center bg-background">
      <div className="mx-auto flex w-full max-w-sm flex-col items-center space-y-4 rounded-xl border bg-card p-8 text-center shadow-lg animate-in fade-in zoom-in duration-500">
        <div className="rounded-full bg-destructive/10 p-3">
          <ShieldAlert className="h-8 w-8 text-destructive" />
        </div>
        <div className="space-y-2">
          <h1 className="text-2xl font-bold tracking-tight">Sessão Expirada</h1>
          <p className="text-sm text-muted-foreground">
            Sua sessão não é mais válida. Estamos redirecionando você para o
            painel
          </p>
        </div>
        <div className="flex items-center space-x-2 text-sm text-muted-foreground pt-4">
          <Loader2 className="h-4 w-4 animate-spin" />
          <span>Redirecionando...</span>
        </div>
      </div>
    </div>
  );
}
