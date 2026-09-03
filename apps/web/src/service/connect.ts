import axios, { AxiosInstance } from "axios";
// import { getApiUrl } from "@/utils/runtime";

export const connectIA = () => {
  const api = axios.create({
    baseURL: "https://painelcomal.duckdns.org/winthor-ia/backend/api/v1",
    // baseURL: "http://localhost:8081/api/v1",
  });

  return api;
};


export const connectDB = (): AxiosInstance => {
  const instance = axios.create({
    // baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1",
    baseURL: "https://painelcomal.duckdns.org/winthor-ia/api/api/v1"
  });

  return instance;
};
