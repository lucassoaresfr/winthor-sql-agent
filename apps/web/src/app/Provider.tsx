"use client";

import React, { createContext, useContext, useState } from "react";
import { ExpiredSession } from "@/components/layout/ExpiredSession";
import { useTokenValidation } from "@/hooks/useTokenValidation";

interface SessionContextType {
  setSessionExpired: () => void;
}

const SessionContext = createContext<SessionContextType>({
  setSessionExpired: () => {},
});

export const useSession = () => useContext(SessionContext);

export function SessionGuard({ children }: { children: React.ReactNode }) {
  const [isExpired, setIsExpired] = useState(false);

  useTokenValidation(1, isExpired, () => setIsExpired(true));

  if (isExpired) {
    return (
      <div className="h-screen w-screen bg-background flex items-center justify-center">
        <ExpiredSession />
      </div>
    );
  }

  return (
    <SessionContext.Provider
      value={{ setSessionExpired: () => setIsExpired(true) }}
    >
      {children}
    </SessionContext.Provider>
  );
}
