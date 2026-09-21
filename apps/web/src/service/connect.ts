import axios, { AxiosInstance } from "axios";

export const PAINEL_PAI_API =
  "https://painelcomal.duckdns.org/api/api/permissoes";

export const connectIA = () => {
  const api = axios.create({
    baseURL: "https://painelcomal.duckdns.org/winthor-ia/backend/api/v1",
    withCredentials: true, // ⚠️ Força o navegador a enviar o cookie HttpOnly automaticamente
  });

  return api;
};

export const connectDB = (): AxiosInstance => {
  const instance = axios.create({
    baseURL: "https://painelcomal.duckdns.org/winthor-ia/api/api/v1",
    withCredentials: true, // ⚠️ Força o navegador a enviar o cookie HttpOnly automaticamente
  });
  return instance;
};
