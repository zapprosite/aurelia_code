import { useEffect, useRef, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";

// Simulated team banter messages (rotate when no real events)
const BANTER: Array<{ from: string; to: string; msg: string }> = [
  { from: "aurelia",    to: "sentinel",    msg: "Status do homelab?" },
  { from: "sentinel",   to: "aurelia",     msg: "Todos containers UP. Load 12%. Sem anomalias." },
  { from: "aurelia",    to: "cronus",      msg: "Crons do dia executados?" },
  { from: "cronus",     to: "aurelia",     msg: "5 jobs agendados. Próximo: sentinel-watchdog em 3min." },
  { from: "aurelia",    to: "gemma",       msg: "Modelos Ollama disponíveis?" },
  { from: "gemma",      to: "aurelia",     msg: "llama3.1:8b ✓  gemma3:12b ✓  VRAM livre: 78%." },
  { from: "openrouter", to: "aurelia",     msg: "MiniMax disponível. Latência 420ms." },
  { from: "sentinel",   to: "cronus",      msg: "GPU 45°C · tudo nominal." },
  { from: "cronus",     to: "sentinel",    msg: "watchdog executado · 0 alertas." },
  { from: "gemma",      to: "openrouter",  msg: "Fallback disponível se necessário?" },
  { from: "openrouter", to: "gemma",       msg: "Afirmativo. Rota de fallback ativa." },
  { from: "aurelia",    to: "all",         msg: "Sistema estável. Continuando monitoramento." },
  { from: "sentinel",   to: "aurelia",     msg: "Disco /srv: 68% livre. ZFS pool saudável." },
  { from: "cronus",     to: "aurelia",     msg: "memory-sync concluído. Qdrant sincronizado." },
  { from: "aurelia",    to: "sentinel",    msg: "Verificar containers parados?" },
  { from: "sentinel",   to: "aurelia",     msg: "Nenhum container parado. 14 rodando." },
];

const AGENT_COLOR: Record<string, string> = {
  aurelia:    "text-purple-400",
  sentinel:   "text-cyan-400",
  cronus:     "text-yellow-400",
  gemma:      "text-green-400",
  openrouter: "text-blue-400",
  all:        "text-white/50",
  system:     "text-white/30",
};

type Msg = {
  id: string;
  from: string;
  to: string;
  msg: string;
  ts: string;
};

function now(): string {
  return new Date().toLocaleTimeString("pt-BR", { hour: "2-digit", minute: "2-digit", second: "2-digit" });
}

function colorOf(name: string) {
  return AGENT_COLOR[name.toLowerCase()] || "text-white/50";
}

export function AgentComms() {
  const [msgs, setMsgs] = useState<Msg[]>([]);
  const banterIdx = useRef(0);
  const bottomRef = useRef<HTMLDivElement>(null);

  const push = (m: Omit<Msg, "id" | "ts">) => {
    setMsgs((prev) => [
      ...prev.slice(-49),
      { ...m, id: Math.random().toString(36).slice(2), ts: now() },
    ]);
  };

  // Rotate banter every 4-7s
  useEffect(() => {
    const fire = () => {
      const b = BANTER[banterIdx.current % BANTER.length];
      banterIdx.current++;
      push(b);
    };
    fire(); // first message immediately
    const iv = setInterval(fire, 4500 + Math.random() * 2500);
    return () => clearInterval(iv);
  }, []);

  // SSE events → agent comms
  useEffect(() => {
    const es = new EventSource("/api/events");
    es.onmessage = (event) => {
      try {
        const d = JSON.parse(event.data);
        if (d.type === "agent_tool" && d.agent) {
          push({ from: d.agent.toLowerCase(), to: "aurelia", msg: `${d.action}` });
        } else if (d.type === "agent_handoff" && d.agent) {
          push({ from: d.agent.toLowerCase(), to: "all", msg: `handoff: ${d.action}` });
        }
      } catch { /* ignore */ }
    };
    return () => es.close();
  }, []);

  // Auto-scroll to bottom
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [msgs]);

  return (
    <div
      className="font-mono text-[10px] rounded-lg border border-white/10 overflow-hidden"
      style={{ background: "rgba(0,0,0,0.5)" }}
    >
      {/* title bar */}
      <div className="flex items-center gap-2 px-3 py-1.5 border-b border-white/10 bg-white/5">
        <span className="text-[9px] text-white/30 uppercase tracking-widest">team radio</span>
        <span className="ml-auto flex items-center gap-1">
          <span className="relative flex h-1.5 w-1.5">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75" />
            <span className="relative inline-flex rounded-full h-1.5 w-1.5 bg-green-500" />
          </span>
          <span className="text-green-400/70 text-[9px]">live</span>
        </span>
      </div>

      {/* messages */}
      <div className="h-44 overflow-y-auto px-3 py-2 space-y-1 scrollbar-thin scrollbar-thumb-white/10">
        <AnimatePresence initial={false}>
          {msgs.map((m) => (
            <motion.div
              key={m.id}
              initial={{ opacity: 0, x: -6 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.2 }}
              className="flex items-start gap-1.5 leading-snug"
            >
              <span className="text-white/20 shrink-0">{m.ts}</span>
              <span className={`shrink-0 font-bold ${colorOf(m.from)}`}>{m.from}</span>
              <span className="text-white/20 shrink-0">→</span>
              <span className={`shrink-0 ${colorOf(m.to)}`}>{m.to}</span>
              <span className="text-white/20 shrink-0">:</span>
              <span className="text-white/60 break-all">{m.msg}</span>
            </motion.div>
          ))}
        </AnimatePresence>
        <div ref={bottomRef} />
      </div>
    </div>
  );
}
