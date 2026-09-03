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
import Image from "next/image";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const currentUserId = "user";

  return (
    <html
      lang="pt-BR"
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased overflow-hidden`}
    >
      <body className="h-full w-full bg-background overflow-hidden flex">
        <TooltipProvider>
          <SidebarProvider defaultOpen={true}>
            {/* 1. Definir a variant como "sidebar" remove a estrutura estilo cartão floating */}
            <AppSidebar userId={currentUserId}/>

            {/* 2. Removemos bordas arredondadas, sombras e margens do SidebarInset */}
            <SidebarInset className="flex flex-col flex-1 h-screen overflow-hidden m-0! rounded-none! shadow-none! border-l">
              <header className="flex h-14 shrink-0 items-center gap-5 border-b px-4 bg-background w-full">
                <SidebarTrigger className="size-4"/>

                {/* Substituição do h1 pela Imagem */}
                <div className="relative flex items-center h-6">
                  <Image
                    src="/LOGO-COLORIDA.png" // Altere para o caminho da sua imagem na pasta /public
                    alt="Winthor Agent AI"
                    width={800} // Ajuste a largura conforme o tamanho do seu logo
                    height={30} // Ajuste a altura proporcionalmente
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
