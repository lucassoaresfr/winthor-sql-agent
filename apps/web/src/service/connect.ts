import axios from "axios";
// import { getApiUrl } from "@/utils/runtime";

export const connectDB = () => {
  const api = axios.create({
//     baseURL: "https://comalext.duckdns.org/check-motorista-api/api/v1",
    baseURL: "http://localhost:8081/api/v1",
  });

  return api;
};
