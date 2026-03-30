import * as React from "react";
import Markdown from "react-markdown";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { oneDark } from "react-syntax-highlighter/dist/esm/styles/prism";
import { Send, Bot, User, Loader2, Sparkles, RotateCcw } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

interface ChatMessage {
  id: string;
  role: "user" | "assistant";
  content: string;
  streaming?: boolean;
  tokens?: { input: number; output: number };
}

export function ChatTab() {
  const [messages, setMessages] = React.useState<ChatMessage[]>([]);
  const [input, setInput] = React.useState("");
  const [isStreaming, setIsStreaming] = React.useState(false);
  const scrollRef = React.useRef<HTMLDivElement>(null);
  const inputRef = React.useRef<HTMLTextAreaElement>(null);

  // Load history on mount
  React.useEffect(() => {
    fetch("/api/chat/history")
      .then((r) => r.json())
      .then((data) => {
        if (Array.isArray(data)) {
          setMessages(
            data.map((m: { role: string; content: string }, i: number) => ({
              id: `hist-${i}`,
              role: m.role as "user" | "assistant",
              content: m.content,
            }))
          );
        }
      })
      .catch(() => {});
  }, []);

  // Auto-scroll to bottom
  React.useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [messages]);

  const sendMessage = async () => {
    const text = input.trim();
    if (!text || isStreaming) return;

    const userMsg: ChatMessage = {
      id: `user-${Date.now()}`,
      role: "user",
      content: text,
    };
    const assistantId = `assistant-${Date.now()}`;
    const assistantMsg: ChatMessage = {
      id: assistantId,
      role: "assistant",
      content: "",
      streaming: true,
    };

    setMessages((prev) => [...prev, userMsg, assistantMsg]);
    setInput("");
    setIsStreaming(true);

    try {
      const response = await fetch("/api/chat", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ message: text }),
      });

      if (!response.ok || !response.body) {
        throw new Error("Failed to connect");
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let buffer = "";

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split("\n");
        buffer = lines.pop() || "";

        for (const line of lines) {
          if (!line.startsWith("data: ")) continue;
          try {
            const token = JSON.parse(line.slice(6));

            if (token.error) {
              setMessages((prev) =>
                prev.map((m) =>
                  m.id === assistantId
                    ? { ...m, content: m.content || token.error, streaming: false }
                    : m
                )
              );
              setIsStreaming(false);
              return;
            }

            if (token.content) {
              setMessages((prev) =>
                prev.map((m) =>
                  m.id === assistantId
                    ? { ...m, content: m.content + token.content }
                    : m
                )
              );
            }

            if (token.done) {
              setMessages((prev) =>
                prev.map((m) =>
                  m.id === assistantId
                    ? {
                        ...m,
                        streaming: false,
                        tokens:
                          token.input_tokens || token.output_tokens
                            ? { input: token.input_tokens, output: token.output_tokens }
                            : undefined,
                      }
                    : m
                )
              );
            }
          } catch {
            // skip malformed lines
          }
        }
      }
    } catch (err) {
      setMessages((prev) =>
        prev.map((m) =>
          m.id === assistantId
            ? { ...m, content: "Erro de conexao. Verifique se o servidor esta rodando.", streaming: false }
            : m
        )
      );
    } finally {
      setIsStreaming(false);
      inputRef.current?.focus();
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  const clearChat = () => {
    setMessages([]);
  };

  return (
    <div className="flex flex-col h-[calc(100vh-8rem)] max-w-4xl mx-auto">
      {/* Messages area */}
      <div ref={scrollRef} className="flex-1 overflow-y-auto space-y-1 pb-4 pr-2 scroll-smooth">
        {messages.length === 0 && (
          <div className="flex flex-col items-center justify-center h-full text-center">
            <div className="w-16 h-16 rounded-2xl bg-primary/10 border border-primary/20 flex items-center justify-center mb-6">
              <Sparkles className="w-8 h-8 text-primary" />
            </div>
            <h2 className="text-xl font-bold text-white/80 mb-2">Aurelia Command Center</h2>
            <p className="text-sm text-white/30 max-w-md">
              Converse diretamente com a Aurelia. Mesmo pipeline do Telegram, com streaming em tempo real.
            </p>
            <div className="flex gap-2 mt-6">
              {["Status do homelab", "Listar containers", "Resumo do dia"].map((suggestion) => (
                <button
                  key={suggestion}
                  onClick={() => {
                    setInput(suggestion);
                    inputRef.current?.focus();
                  }}
                  className="px-3 py-1.5 rounded-lg bg-white/5 border border-white/10 text-xs text-white/40 hover:text-white/70 hover:border-white/20 transition-colors"
                >
                  {suggestion}
                </button>
              ))}
            </div>
          </div>
        )}

        <AnimatePresence initial={false}>
          {messages.map((msg) => (
            <motion.div
              key={msg.id}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.2 }}
              className={`flex gap-3 py-4 px-2 ${
                msg.role === "user" ? "" : ""
              }`}
            >
              {/* Avatar */}
              <div
                className={`w-8 h-8 rounded-lg flex items-center justify-center shrink-0 ${
                  msg.role === "user"
                    ? "bg-white/10 border border-white/10"
                    : "bg-primary/10 border border-primary/20"
                }`}
              >
                {msg.role === "user" ? (
                  <User className="w-4 h-4 text-white/50" />
                ) : (
                  <Bot className="w-4 h-4 text-primary" />
                )}
              </div>

              {/* Content */}
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <span className="text-xs font-bold text-white/50 uppercase tracking-wider">
                    {msg.role === "user" ? "Voce" : "Aurelia"}
                  </span>
                  {msg.streaming && (
                    <Loader2 className="w-3 h-3 text-primary animate-spin" />
                  )}
                  {msg.tokens && (
                    <span className="text-[10px] text-white/20 font-mono">
                      {msg.tokens.input}in/{msg.tokens.output}out
                    </span>
                  )}
                </div>

                {msg.role === "user" ? (
                  <p className="text-sm text-white/70 whitespace-pre-wrap">{msg.content}</p>
                ) : (
                  <div className="prose prose-invert prose-sm max-w-none text-white/70 [&_p]:mb-2 [&_p]:leading-relaxed [&_ul]:mb-2 [&_ol]:mb-2 [&_li]:mb-0.5 [&_pre]:bg-[#1e1e2e] [&_pre]:border [&_pre]:border-white/10 [&_pre]:rounded-lg">
                    <Markdown
                      components={{
                        code({ className, children, ...props }) {
                          const match = /language-(\w+)/.exec(className || "");
                          const inline = !match;
                          return inline ? (
                            <code
                              className="px-1.5 py-0.5 rounded bg-white/10 text-primary/90 text-xs font-mono"
                              {...props}
                            >
                              {children}
                            </code>
                          ) : (
                            <SyntaxHighlighter
                              style={oneDark}
                              language={match[1]}
                              PreTag="div"
                              customStyle={{
                                margin: 0,
                                borderRadius: "0.5rem",
                                fontSize: "0.75rem",
                              }}
                            >
                              {String(children).replace(/\n$/, "")}
                            </SyntaxHighlighter>
                          );
                        },
                      }}
                    >
                      {msg.content || (msg.streaming ? "" : "...")}
                    </Markdown>
                    {msg.streaming && msg.content && (
                      <span className="inline-block w-1.5 h-4 bg-primary/60 animate-pulse rounded-sm ml-0.5 align-text-bottom" />
                    )}
                  </div>
                )}
              </div>
            </motion.div>
          ))}
        </AnimatePresence>
      </div>

      {/* Input bar */}
      <div className="border-t border-white/10 pt-4 pb-2">
        <div className="flex items-end gap-3">
          <div className="flex-1 relative">
            <textarea
              ref={inputRef}
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Envie uma mensagem para a Aurelia..."
              rows={1}
              className="w-full resize-none rounded-xl bg-white/5 border border-white/10 px-4 py-3 pr-12 text-sm text-white/80 placeholder-white/20 focus:outline-none focus:border-primary/40 focus:ring-1 focus:ring-primary/20 transition-all"
              style={{ minHeight: "48px", maxHeight: "120px" }}
              onInput={(e) => {
                const target = e.target as HTMLTextAreaElement;
                target.style.height = "auto";
                target.style.height = Math.min(target.scrollHeight, 120) + "px";
              }}
            />
          </div>
          <div className="flex gap-2">
            {messages.length > 0 && (
              <button
                onClick={clearChat}
                className="p-3 rounded-xl bg-white/5 border border-white/10 text-white/30 hover:text-white/60 hover:bg-white/10 transition-colors"
                title="Limpar conversa"
              >
                <RotateCcw className="w-4 h-4" />
              </button>
            )}
            <button
              onClick={sendMessage}
              disabled={isStreaming || !input.trim()}
              className="p-3 rounded-xl bg-primary/20 border border-primary/30 text-primary hover:bg-primary/30 transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
            >
              {isStreaming ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : (
                <Send className="w-4 h-4" />
              )}
            </button>
          </div>
        </div>
        <div className="flex items-center justify-between mt-2 px-1">
          <span className="text-[10px] text-white/15 font-mono">
            Enter para enviar / Shift+Enter para nova linha
          </span>
          <span className="text-[10px] text-white/15 font-mono uppercase tracking-widest">
            Jarvis Command Center
          </span>
        </div>
      </div>
    </div>
  );
}
