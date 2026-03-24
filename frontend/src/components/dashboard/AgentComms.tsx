import { useEffect, useRef, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";

// Simulated team banter messages (rotate when no real events)
const BANTER: Array<{ from: string; to: string; msg: string }> = [
  { from: "aurelia",    to: "sentinel",    msg: "Status do homelab?" },
  { from: "sentinel",   to: "aurelia",     msg: "Todos containers UP. Load 8%. ZFS saudável." },
  { from: "aurelia",    to: "cronus",      msg: "Jobs de manutenção agendados?" },
  { from: "cronus",     to: "aurelia",     msg: "4 jobs ativos. Sentinel-watchdog em T-2min." },
  { from: "aurelia",    to: "gemma",       msg: "Estado do motor Tier 0?" },
  { from: "gemma",      to: "aurelia",     msg: "gemma3:12b (Resident) ✓  VRAM: 8.2GB livre." },
  { from: "openrouter", to: "aurelia",     msg: "Tier 1 (DeepSeek) e Tier 2 (MiniMax) operacionais." },
  { from: "sentinel",   to: "cronus",      msg: "GPU 42°C · drivers NVIDIA v550.67." },
  { from: "cronus",     to: "sentinel",    msg: "Watchdog executado · 0 alertas · Logs limpos." },
  { from: "gemma",      to: "openrouter",  msg: "Sincronia de contexto via Qdrant ativa." },
  { from: "openrouter", to: "gemma",       msg: "Confirmado. Handoff de memória habilitado." },
  { from: "aurelia",    to: "all",         msg: "Governança Industrial 2026: Sistema Nominal." },
  { from: "sentinel",   to: "aurelia",     msg: "Disco /srv: 72% livre. I/O estável." },
  { from: "cronus",     to: "aurelia",     msg: "Backup de banco concluído em 42ms." },
  { from: "aurelia",    to: "sentinel",    msg: "Consistência de rede Cloudflare?" },
  { from: "sentinel",   to: "aurelia",     msg: "Tunnel Estável. 0 pacotes perdidos em 24h." },
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
        } else if (d.type === "agent_comms" && d.agent) {
          // S-25: real cron/health events from backend
          const payload = d.payload || {};
          const jobId = typeof payload === "object" ? (payload.job || "") : "";
          const status = typeof payload === "object" ? (payload.status || "ok") : "ok";
          const shortJob = jobId.length > 8 ? jobId.slice(0, 8) + "…" : jobId;
          push({
            from: d.agent.toLowerCase(),
            to: "aurelia",
            msg: `★ ${d.action}${shortJob ? ` [${shortJob}]` : ""} — ${status}`,
          });
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
