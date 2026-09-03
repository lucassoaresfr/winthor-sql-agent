import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { TooltipProvider } from "@/components/ui/tooltip";
import { Toaster } from "@/components/ui/toast";
import {
  SidebarProvider,
  SidebarInset,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/layout/sidebar/app-sidebar";
import { ExpiredSession } from "@/components/layout/ExpiredSession";
import Image from "next/image";
import logocomal from "@/../public/LOGO-COLORIDA.png";
import { getTokenPayload } from "@/service/TokenValid";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Winthor Agent AI",
  description: "Assistente virtual do ERP Winthor",
  icons: {
    icon: "/logo.svg",
  },
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const token = await getTokenPayload();

  if (!token) {
    return (
      <html
        lang="pt-BR"
        className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
      >
        <body className="h-full w-full bg-background flex items-center justify-center">
          <ExpiredSession />
        </body>
      </html>
    );
  }

  // Define o userId como o NOME do usuário (com fallback para login/idPg caso venha vazio)
  const currentUserName = token.nome || token.usuario || String(token.idPg);
  const currentUserLogin = token.usuario;

  return (
    <html
      lang="pt-BR"
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased overflow-hidden`}
    >
      <body className="h-full w-full bg-background overflow-hidden flex">
        <TooltipProvider>
          <SidebarProvider defaultOpen={true}>
            {/* Passamos o NOME em userId e opcionalmente o login/matrícula */}
            <AppSidebar userId={currentUserName} userLogin={currentUserLogin} />

            <SidebarInset className="flex flex-col flex-1 h-screen overflow-hidden m-0! rounded-none! shadow-none! border-l">
              <header className="flex h-14 shrink-0 items-center gap-5 border-b px-4 bg-background w-full">
                <SidebarTrigger className="size-4" />

                <div className="relative flex items-center h-6">
                  <Image
                    src={logocomal}
                    alt="Winthor Agent AI"
                    width={800}
                    height={30}
                    className="object-contain h-full w-auto"
                    priority
                  />
                </div>
              </header>

              <main className="flex-1 overflow-hidden relative">
                {children}
              </main>
            </SidebarInset>
          </SidebarProvider>
        </TooltipProvider>
        <Toaster />
      </body>
    </html>
  );
}
