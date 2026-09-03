"use client";

import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { cn } from "@/lib/utils";

interface MarkdownRendererProps {
  content: string;
  className?: string;
}

export function MarkdownRenderer({
  content,
  className,
}: MarkdownRendererProps) {
  const formattedContent = content
    .replace(/\\n/g, "\n")
    .replace(/^[ \t]{4}/gm, "")
    .replace(/([^\n])\s*\*\s+\*\*/g, "$1\n* **")
    .trim();

  return (
    <div
      className={cn(
        "text-sm text-stone-800 leading-relaxed wrap-break-word",
        className,
      )}
    >
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        components={{
          // Renderização estilizada de Tabelas
          table: ({ children }) => (
            <div className="my-3 overflow-x-auto rounded-lg border border-stone-200 shadow-sm">
              <table className="w-full text-left border-collapse text-xs sm:text-sm">
                {children}
              </table>
            </div>
          ),
          thead: ({ children }) => (
            <thead className="bg-stone-100 border-b border-stone-200 text-stone-700 font-semibold">
              {children}
            </thead>
          ),
          tbody: ({ children }) => (
            <tbody className="divide-y divide-stone-200 bg-white">
              {children}
            </tbody>
          ),
          tr: ({ children }) => (
            <tr className="hover:bg-stone-50/50 transition-colors">
              {children}
            </tr>
          ),
          th: ({ children }) => (
            <th className="px-3.5 py-2.5 font-medium">{children}</th>
          ),
          td: ({ children }) => (
            <td className="px-3.5 py-2.5 text-stone-700">{children}</td>
          ),

          // Renderização de parágrafos
          p: ({ children }) => <div className="mb-2 last:mb-0">{children}</div>,
          strong: ({ children }) => (
            <strong className="font-semibold text-stone-900">{children}</strong>
          ),
          em: ({ children }) => (
            <em className="italic text-stone-700">{children}</em>
          ),

          // Listas
          ul: ({ children }) => (
            <ul className="list-disc space-y-1 my-2 pl-5 text-stone-700">
              {children}
            </ul>
          ),
          ol: ({ children }) => (
            <ol className="list-decimal space-y-1 my-2 pl-5 text-stone-700">
              {children}
            </ol>
          ),
          li: ({ children }) => (
            <li className="pl-1 text-stone-700">{children}</li>
          ),

          // Bloco de Código (<pre>)
          pre({ children }) {
            return (
              <pre className="my-3 p-3.5 bg-stone-900 text-stone-100 rounded-xl font-mono text-xs overflow-x-auto shadow-inner">
                {children}
              </pre>
            );
          },

          // Tag <code> (Inline vs Block)
          code({ node, inline, className, children, ...props }: any) {
            // Se for inline (ex: `codigo`), renderiza tag discreta
            if (inline) {
              return (
                <code
                  className={cn(
                    "px-1.5 py-0.5 bg-stone-100 text-stone-800 rounded border border-stone-200 font-mono text-xs",
                    className,
                  )}
                  {...props}
                >
                  {children}
                </code>
              );
            }

            // Se for dentro de um bloco <pre>, limpa classes extras para evitar caixas duplicadas
            return (
              <code className={cn("font-mono text-xs", className)} {...props}>
                {children}
              </code>
            );
          },
        }}
      >
        {formattedContent}
      </ReactMarkdown>
    </div>
  );
}
