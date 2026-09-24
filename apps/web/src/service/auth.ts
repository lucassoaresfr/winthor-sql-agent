import axios from "axios";

const LOGOUT_URL =
  "/api/api/permissoes/auth/logout";
const LOGIN_URL = "/";

// Flag global em memória para impedir requisições/redirecionamentos simultâneos
let isLoggingOut = false;

export async function performLogout() {
  if (isLoggingOut) return;
  isLoggingOut = true;

  try {
    // 1. Invoca a API pai para destruir o cookie HttpOnly no backend
    await axios.post(LOGOUT_URL, {}, { withCredentials: true });
  } catch (error) {
    console.error("[Auth] Erro ao executar logout no servidor:", error);
  } finally {
    // 3. Redireciona e substitui o histórico para a rota de login "/"
    if (typeof window !== "undefined") {
      window.location.replace(LOGIN_URL);
    }
  }
}