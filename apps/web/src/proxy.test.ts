import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const PAINEL_PAI_API =
  process.env.PAINEL_PAI_API_URL ||
  "https://painelcomal.duckdns.org/painel-test/api/api/permissoes";

export async function proxy(request: NextRequest) {
  const cookieHeader = request.headers.get("cookie");
  const token = request.cookies.get("token")?.value;

  let isExpired = false;

  if (!token) {
    isExpired = true;
  } else {
    try {
      const cleanToken = token
        .replace(/^["']|["']$/g, "")
        .replace(/^Bearer\s+/i, "")
        .trim();

      const response = await fetch(`${PAINEL_PAI_API}/auth/validate`, {
        method: "GET",
        headers: {
          Authorization: `Bearer ${cleanToken}`,
          ...(cookieHeader ? { cookie: cookieHeader } : {}),
        },
        cache: "no-store",
      });

      if (!response.ok) {
        isExpired = true;
      }
    } catch (error) {
      console.error("[Projeto Filho] Erro na validação do token:", error);
      isExpired = true;
    }
  }

  if (isExpired) {
    // 1. Faz a requisição POST ao endpoint de logout do Express Pai
    // Isso garante que o backend invalide o token e injete o Set-Cookie para expirar o HttpOnly
    if (cookieHeader || token) {
      try {
        await fetch(`${PAINEL_PAI_API}/auth/logout`, {
          method: "POST",
          headers: {
            ...(cookieHeader ? { cookie: cookieHeader } : {}),
            ...(token ? { Authorization: `Bearer ${token}` } : {}),
          },
          cache: "no-store",
        });
      } catch (logoutError) {
        console.error(
          "[Projeto Filho] Erro ao chamar a rota de logout:",
          logoutError,
        );
      }
    }

    // 2. Prepara os headers de resposta para garantir que o Server Component leia a sessão expirada
    const requestHeaders = new Headers(request.headers);
    requestHeaders.set("x-session-expired", "true");

    const response = NextResponse.next({
      request: {
        headers: requestHeaders,
      },
    });

    // 3. Define o header na resposta visível ao RootLayout/Client
    response.headers.set("x-session-expired", "true");

    // 4. Garante a expiração do cookie no nível da resposta do Next.js
    response.cookies.set("token", "", {
      path: "/",
      domain: ".duckdns.org",
      expires: new Date(0),
      httpOnly: true,
    });

    return response;
  }

  return NextResponse.next();
}

// Matcher garantindo a interceptação em todas as páginas, ignorando recursos estáticos e chamadas de API
export const config = {
  matcher: [
    "/((?!api|_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)",
  ],
};
