"use client";

import { useEffect, useRef } from "react";
import axios from "axios";

const VALIDATE_URL =
  "https://painelcomal.duckdns.org/painel-test/api/api/permissoes/auth/validate";

export function useTokenValidation(
  intervalInMinutes: number = 1,
  isAlreadyExpired: boolean = false,
  onExpire?: () => void,
) {
  const isExpiredRef = useRef(isAlreadyExpired);

  useEffect(() => {
    isExpiredRef.current = isAlreadyExpired;
  }, [isAlreadyExpired]);

  useEffect(() => {
    if (isAlreadyExpired || isExpiredRef.current) return;

    const checkToken = async () => {
      if (isExpiredRef.current) return;

      try {
        await axios.get(VALIDATE_URL, {
          withCredentials: true,
        });
      } catch (error: any) {
        if (
          error.response?.status === 401 ||
          error.response?.status === 403 ||
          !error.response
        ) {
          console.warn("[Session] Token expirado detectado pelo polling.");
          isExpiredRef.current = true;

          // Executa APENAS a troca de estado no React
          if (onExpire) {
            onExpire();
          }
        }
      }
    };

    checkToken();

    const intervalMs = intervalInMinutes * 60 * 1000;
    const timer = setInterval(checkToken, intervalMs);

    return () => clearInterval(timer);
  }, [intervalInMinutes, isAlreadyExpired, onExpire]);
}
