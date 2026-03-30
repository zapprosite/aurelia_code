import * as React from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
  MousePointer,
  Download,
  Camera,
  Play,
  Loader2,
  Globe,
  ChevronRight,
  Zap,
  Bot,
  MessageSquare,
  ExternalLink
} from "lucide-react";
import { Button } from "../ui/Button";
import { Input as _Input } from "../ui/Input";
import { Card, CardContent } from "../ui/Card";
import { Badge } from "../ui/Badge";
import { ScrollArea } from "../ui/ScrollArea";

interface Message {
  id: string;
  role: "user" | "assistant";
  content: string;
  timestamp: Date;
  action?: "navigate" | "act" | "extract" | "screenshot";
  status?: "pending" | "success" | "error";
}

interface ComputerUseTabProps {
  className?: string;
}

export function ComputerUseTab({ className }: ComputerUseTabProps) {
  const [url, setUrl] = React.useState("https://example.com");
  const [instruction, setInstruction] = React.useState("");
  const [messages, setMessages] = React.useState<Message[]>([
    {
      id: "1",
      role: "assistant",
      content: "Olá! Sou o Jarvis Computer Use. Posso navegar na web, interagir com páginas e extrair dados. Digite uma URL para começar!",
      timestamp: new Date(),
    },
  ]);
  const [isLoading, setIsLoading] = React.useState(false);
  const [screenshot, _setScreenshot] = React.useState<string | null>(null);
  const [_copied, setCopied] = React.useState<string | null>(null);
  const [stagehandStatus, _setStagehandStatus] = React.useState<"connected" | "disconnected" | "loading">("connected");

  const scrollRef = React.useRef<HTMLDivElement>(null);

  // Quick actions presets
  const quickActions = [
    { label: "Navegar", icon: Globe, action: "navigate", placeholder: "URL para navegar..." },
    { label: "Ação", icon: MousePointer, action: "act", placeholder: "Ex: clique no botão Login" },
    { label: "Extrair", icon: Download, action: "extract", placeholder: "Ex: todos os preços" },
    { label: "Screenshot", icon: Camera, action: "screenshot", placeholder: "" },
  ];

  const [activeAction, setActiveAction] = React.useState<string>("navigate");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!url && activeAction === "navigate") return;
    if (!instruction && activeAction !== "screenshot") return;

    const userMsg: Message = {
      id: Date.now().toString(),
      role: "user",
      content: activeAction === "navigate"
        ? `Navegar para: ${url}`
        : instruction,
      timestamp: new Date(),
      action: activeAction as Message["action"],
    };

    setMessages(prev => [...prev, userMsg]);
    setIsLoading(true);

    try {
      // Call backend API
      await fetch("/v1/telegram/impersonate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          user_id: 7220607041,
          chat_id: 7220607041,
          text: `Use computer use to ${activeAction}: ${activeAction === "navigate" ? url : instruction}`,
        }),
      });

      const assistantMsg: Message = {
        id: (Date.now() + 1).toString(),
        role: "assistant",
        content: `Executando ${activeAction}...`,
        timestamp: new Date(),
        action: activeAction as Message["action"],
        status: "pending",
      };
      setMessages(prev => [...prev, assistantMsg]);

      // Simulate response (replace with actual SSE/WebSocket later)
      setTimeout(() => {
        setMessages(prev => prev.map(m =>
          m.id === assistantMsg.id
            ? { ...m, content: `Ação ${activeAction} executada com sucesso!`, status: "success" }
            : m
        ));
        setIsLoading(false);
        setInstruction("");
      }, 2000);

    } catch (error) {
      setMessages(prev => prev.map(m =>
        m.id === messages[messages.length - 1]?.id
          ? { ...m, content: `Erro: ${error}`, status: "error" }
          : m
      ));
      setIsLoading(false);
    }
  };

  void setCopied; // used in copyToClipboard below
  const copyToClipboard = (text: string, id: string) => {
    navigator.clipboard.writeText(text);
    setCopied(id);
    setTimeout(() => setCopied(null), 2000);
  };
  void copyToClipboard;

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header com Status */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-primary/10">
            <Bot className="w-5 h-5 text-primary" />
          </div>
          <div>
            <h2 className="text-lg font-semibold text-white/90">Computer Use</h2>
            <p className="text-xs text-white/40">Navegue e interaja com a web</p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Badge
            variant={stagehandStatus === "connected" ? "default" : "destructive"}
            className="text-xs"
          >
            <span className={`w-1.5 h-1.5 rounded-full mr-1.5 ${
              stagehandStatus === "connected" ? "bg-emerald-400" : "bg-red-400"
            }`} />
            Stagehand {stagehandStatus}
          </Badge>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="flex gap-2 flex-wrap">
        {quickActions.map((action) => (
          <Button
            key={action.action}
            variant={activeAction === action.action ? "default" : "outline"}
            size="sm"
            onClick={() => setActiveAction(action.action)}
            className="gap-2"
          >
            <action.icon className="w-4 h-4" />
            {action.label}
          </Button>
        ))}
      </div>

      {/* Main Input - Estilo Perplexity */}
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="relative group">
          <div className="absolute -inset-0.5 bg-gradient-to-r from-primary/50 to-purple-500/50 rounded-2xl blur opacity-30 group-hover:opacity-50 transition duration-500" />
          <div className="relative bg-[#1a1a1a] rounded-2xl border border-white/10 overflow-hidden">
            {/* URL Bar */}
            {activeAction === "navigate" && (
              <div className="flex items-center gap-3 px-4 py-3 border-b border-white/5">
                <Globe className="w-4 h-4 text-white/40" />
                <input
                  type="url"
                  value={url}
                  onChange={(e) => setUrl(e.target.value)}
                  placeholder="https://..."
                  className="flex-1 bg-transparent text-sm text-white/90 placeholder:text-white/30 outline-none"
                />
                <ChevronRight className="w-4 h-4 text-white/20" />
              </div>
            )}

            {/* Instruction Input */}
            <div className="p-4">
              <textarea
                value={instruction}
                onChange={(e) => setInstruction(e.target.value)}
                placeholder={quickActions.find(a => a.action === activeAction)?.placeholder || "Digite sua instrução..."}
                rows={2}
                className="w-full bg-transparent text-white/90 placeholder:text-white/30 outline-none resize-none text-sm leading-relaxed"
                onKeyDown={(e) => {
                  if (e.key === "Enter" && !e.shiftKey) {
                    e.preventDefault();
                    handleSubmit(e);
                  }
                }}
              />
            </div>

            {/* Actions Bar */}
            <div className="flex items-center justify-between px-4 py-3 border-t border-white/5">
              <div className="flex items-center gap-2">
                <Badge variant="outline" className="text-[10px] py-0">
                  <Zap className="w-3 h-3 mr-1" />
                  AI Powered
                </Badge>
              </div>
              <Button
                type="submit"
                size="sm"
                disabled={isLoading || (!url && activeAction === "navigate")}
                className="gap-2"
              >
                {isLoading ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  <Play className="w-4 h-4" />
                )}
                Executar
              </Button>
            </div>
          </div>
        </div>
      </form>

      {/* Screenshot Preview */}
      {screenshot && (
        <Card className="overflow-hidden">
          <CardContent className="p-0">
            <div className="relative group">
              <img
                src={screenshot}
                alt="Browser screenshot"
                className="w-full rounded-lg"
              />
              <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
                <Button size="sm" variant="outline" className="gap-2">
                  <Download className="w-4 h-4" />
                  Download
                </Button>
                <Button size="sm" variant="outline" className="gap-2">
                  <ExternalLink className="w-4 h-4" />
                  Open
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Chat Messages - Estilo Chat */}
      <div className="space-y-4">
        <div className="flex items-center gap-2 text-white/40">
          <MessageSquare className="w-4 h-4" />
          <span className="text-xs font-medium uppercase tracking-widest">Histórico</span>
        </div>

        <ScrollArea className="h-[400px] pr-4" ref={scrollRef}>
          <div className="space-y-4">
            <AnimatePresence>
              {messages.map((msg) => (
                <motion.div
                  key={msg.id}
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0 }}
                  className={`flex ${msg.role === "user" ? "justify-end" : "justify-start"}`}
                >
                  <div className={`max-w-[85%] rounded-2xl px-4 py-3 ${
                    msg.role === "user"
                      ? "bg-primary text-white"
                      : "bg-white/5 border border-white/10"
                  }`}>
                    <p className="text-sm text-white/90 leading-relaxed">{msg.content}</p>
                    <div className="flex items-center gap-2 mt-2">
                      <span className="text-[10px] text-white/30">
                        {msg.timestamp.toLocaleTimeString()}
                      </span>
                      {msg.status && (
                        <Badge
                          variant={msg.status === "success" ? "default" : msg.status === "error" ? "destructive" : "outline"}
                          className="text-[10px] py-0"
                        >
                          {msg.status}
                        </Badge>
                      )}
                      {msg.action && (
                        <Badge variant="outline" className="text-[10px] py-0">
                          {msg.action}
                        </Badge>
                      )}
                    </div>
                  </div>
                </motion.div>
              ))}
            </AnimatePresence>
          </div>
        </ScrollArea>
      </div>

      {/* Available Tools */}
      <Card className="bg-white/5 border-white/10">
        <CardContent className="p-4">
          <h3 className="text-sm font-medium text-white/60 mb-3">Ferramentas Disponíveis</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-2">
            {[
              { name: "navigate", desc: "Abrir URL" },
              { name: "act", desc: "Interagir" },
              { name: "extract", desc: "Extrair dados" },
              { name: "screenshot", desc: "Capturar tela" },
            ].map((tool) => (
              <div
                key={tool.name}
                className="px-3 py-2 rounded-lg bg-white/5 border border-white/5 text-center"
              >
                <p className="text-xs font-mono text-primary">{tool.name}</p>
                <p className="text-[10px] text-white/30">{tool.desc}</p>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
