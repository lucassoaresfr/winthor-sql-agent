import type { Metadata, Viewport } from "next";
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
import Image from "next/image";
import logocomal from "@/../public/LOGO-COLORIDA.png";
import { getTokenPayload } from "@/service/TokenValid";
import { SessionGuard } from "./Provider";

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

// Evita o zoom automático incômodo ao focar no campo de digitação em mobiles
export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
  maximumScale: 1,
  userScalable: false,
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const token = await getTokenPayload();

  const currentUserName =
    token?.nome || token?.usuario || String(token?.idPg ?? "");
  const currentUserLogin = token?.usuario || "";

  return (
    <html
      lang="pt-BR"
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased overflow-hidden`}
    >
      <body className="h-dvh w-full bg-background overflow-hidden flex">
        <SessionGuard>
          <TooltipProvider>
            {/* Em telas menores o menu por padrão começa recolhido para dar espaço ao chat */}
            <SidebarProvider defaultOpen={true}>
              <AppSidebar
                userId={currentUserName}
                userLogin={currentUserLogin}
              />

              <SidebarInset className="flex flex-col flex-1 h-dvh overflow-hidden m-0! rounded-none! shadow-none! border-l-0 sm:border-l">
                {/* Header Responsivo */}
                <header className="flex h-14 shrink-0 items-center justify-between sm:justify-start gap-3 sm:gap-5 border-b px-3 sm:px-4 bg-background w-full">
                  <div className="flex items-center gap-3">
                    {/* Botão de abrir/fechar a Sidebar bem visível no mobile */}
                    <SidebarTrigger className="size-8 sm:size-4 p-1 rounded-md border sm:border-none shadow-xs sm:shadow-none" />

                    <div className="relative flex items-center h-5 sm:h-6">
                      <Image
                        src={logocomal}
                        alt="Winthor Agent AI"
                        width={800}
                        height={30}
                        className="object-contain h-full w-auto max-w-30 sm:max-w-none"
                        priority
                      />
                    </div>
                  </div>
                </header>

                <main className="flex-1 overflow-hidden relative h-[calc(100dvh-3.5rem)]">
                  {children}
                </main>
              </SidebarInset>
            </SidebarProvider>
          </TooltipProvider>
          <Toaster />
        </SessionGuard>
      </body>
    </html>
  );
}
