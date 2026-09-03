import { cookies } from "next/headers";
import { jwtDecode } from "jwt-decode";

export interface TokenPayload {
  nome: string;
  usuario: string;
  codsetor: number | null;
  codsetorPg: number | null;
  email: string;
  idPg: number;
}

export async function getTokenPayload(): Promise<TokenPayload | null> {
  const cookieStore = await cookies();
  const tokenCookie = cookieStore.get("token");

  if (!tokenCookie || !tokenCookie?.value) {
    return null;
  }

  try {
    return jwtDecode<TokenPayload>(tokenCookie.value);
  } catch (error) {
    console.error("Token inválido ou corrompido no servidor:", error);
    return null;
  }
}
