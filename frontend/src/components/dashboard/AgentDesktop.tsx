import { useEffect, useState } from "react";
import * as Icons from "lucide-react";
import { motion } from "framer-motion";
import { type SquadAgent } from "./SquadGrid";

// Fixed desk layout for grid view
const DESK_LAYOUT: Record<string, { col: number; row: number }> = {
  aurelia:    { col: 3, row: 1 },
  sentinel:   { col: 1, row: 3 },
  cronus:     { col: 5, row: 3 },
  gemma:      { col: 2, row: 3 },
  openrouter: { col: 4, row: 3 },
};

const FREE_SLOTS: Array<{ col: number; row: number }> = [
  { col: 1, row: 2 },
  { col: 5, row: 2 },
  { col: 1, row: 4 },
  { col: 5, row: 4 },
  { col: 2, row: 2 },
  { col: 4, row: 2 },
];

const getIcon = (iconName: string) => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const Icon = (Icons as Record<string, any>)[iconName] || Icons.CircleDashed;
  return Icon;
};

function agentSlot(agent: SquadAgent, usedSlots: Set<string>): { col: number; row: number } {
  const key = agent.id.toLowerCase();
  if (DESK_LAYOUT[key]) return DESK_LAYOUT[key];
  for (const slot of FREE_SLOTS) {
    const slotKey = `${slot.col}-${slot.row}`;
    if (!usedSlots.has(slotKey)) {
      usedSlots.add(slotKey);
      return slot;
    }
  }
  return { col: 6, row: 4 };
}

function asciiBar(load: number, width = 10): string {
  const filled = Math.round((load / 100) * width);
  return "[" + "█".repeat(filled) + "░".repeat(width - filled) + "]";
}

function statusDot(status: string): string {
  switch (status) {
    case "online":  return "●";
    case "busy":    return "⟳";
    case "offline": return "◌";
    default:        return "·";
  }
}

function borderColor(color: string): string {
  const map: Record<string, string> = {
    "text-purple-400": "border-purple-500/60",
    "text-cyan-400":   "border-cyan-500/60",
    "text-yellow-400": "border-yellow-500/60",
    "text-green-400":  "border-green-500/60",
    "text-blue-400":   "border-blue-500/60",
  };
  return map[color] || "border-white/20";
}

function titleColor(color: string): string {
  const map: Record<string, string> = {
    "text-purple-400": "bg-purple-900/40 text-purple-300",
    "text-cyan-400":   "bg-cyan-900/40 text-cyan-300",
    "text-yellow-400": "bg-yellow-900/40 text-yellow-300",
    "text-green-400":  "bg-green-900/40 text-green-300",
    "text-blue-400":   "bg-blue-900/40 text-blue-300",
  };
  return map[color] || "bg-white/10 text-white/60";
}

function statusTextColor(status: string): string {
  switch (status) {
    case "online":  return "text-green-400";
    case "busy":    return "text-yellow-400";
    case "offline": return "text-white/30";
    default:        return "text-white/30";
  }
}

// ── TmuxCard ─────────────────────────────────────────────────────────────────

function TmuxCard({ agent, delay = 0 }: { agent: SquadAgent; delay?: number }) {
  const Icon = getIcon(agent.icon);
  const bar = asciiBar(agent.load);

  return (
    <motion.div
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: agent.status === "offline" ? 0.4 : 1, y: 0 }}
      transition={{ delay, duration: 0.25 }}
      className={`font-mono text-xs border rounded overflow-hidden ${borderColor(agent.color)}`}
      style={{ background: "rgba(0,0,0,0.55)" }}
    >
      {/* tmux title bar */}
      <div className={`flex items-center gap-2 px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider ${titleColor(agent.color)}`}>
        <Icon className="w-3 h-3 flex-shrink-0" />
        <span>[{agent.id.toUpperCase()}]</span>
        <span className="ml-auto opacity-60">─</span>
      </div>

      {/* body */}
      <div className="px-3 py-2 space-y-1">
        <div className="text-white/80 font-semibold text-[11px]">{agent.name}</div>
        <div className="text-white/35 text-[9px] leading-tight">{agent.role}</div>

        <div className="pt-1 space-y-0.5">
          <div className={`${statusTextColor(agent.status)} text-[10px]`}>
            {statusDot(agent.status)} {agent.status}
          </div>
          <div className="text-white/40 text-[9px] tracking-tighter">
            {bar} {agent.load}%
          </div>
        </div>
      </div>
    </motion.div>
  );
}

// ── TrelloBoard ───────────────────────────────────────────────────────────────

const COLUMNS = [
  { key: "online",  label: "● ONLINE",  color: "text-green-400",  border: "border-green-500/30" },
  { key: "busy",    label: "⟳ BUSY",    color: "text-yellow-400", border: "border-yellow-500/30" },
  { key: "offline", label: "◌ OFFLINE", color: "text-white/30",   border: "border-white/10" },
];

function TrelloBoard({ agents }: { agents: SquadAgent[] }) {
  return (
    <div className="grid grid-cols-3 gap-3 font-mono">
      {COLUMNS.map((col) => {
        const colAgents = agents.filter((a) => a.status === col.key);
        return (
          <div key={col.key} className={`border rounded-lg p-2 ${col.border}`} style={{ background: "rgba(0,0,0,0.3)" }}>
            {/* column header */}
            <div className={`text-[10px] font-bold uppercase tracking-widest mb-2 px-1 ${col.color}`}>
              {col.label}
              <span className="ml-1 opacity-40">({colAgents.length})</span>
            </div>
            <div className="space-y-2">
              {colAgents.length === 0 ? (
                <div className="text-white/15 text-[9px] px-1 py-2 text-center">— empty —</div>
              ) : (
                colAgents.map((agent, i) => (
                  <TmuxCard key={agent.id} agent={agent} delay={i * 0.06} />
                ))
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}

// ── DesktopGrid (original 2D layout) ─────────────────────────────────────────

function DesktopGrid({ agents }: { agents: SquadAgent[] }) {
  const usedSlots = new Set<string>(
    Object.values(DESK_LAYOUT).map((s) => `${s.col}-${s.row}`)
  );
  const agentSlots = agents.map((agent) => ({
    agent,
    slot: agentSlot(agent, usedSlots),
  }));

  return (
    <div className="relative w-full" style={{ minHeight: "420px" }}>
      <div
        className="absolute inset-0 pointer-events-none rounded-xl opacity-20"
        style={{
          backgroundImage:
            "repeating-linear-gradient(0deg, transparent, transparent 79px, rgba(255,255,255,0.04) 80px), " +
            "repeating-linear-gradient(90deg, transparent, transparent 79px, rgba(255,255,255,0.04) 80px)",
          backgroundSize: "80px 80px",
        }}
      />
      <div
        className="relative grid gap-3 p-3"
        style={{
          gridTemplateColumns: "repeat(6, 1fr)",
          gridTemplateRows: "repeat(4, auto)",
        }}
      >
        {agentSlots.map(({ agent, slot }, i) => (
          <motion.div
            key={agent.id}
            initial={{ opacity: 0, scale: 0.85 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: i * 0.08 }}
            style={{ gridColumn: slot.col, gridRow: slot.row }}
          >
            <TmuxCard agent={agent} />
          </motion.div>
        ))}
      </div>
    </div>
  );
}

// ── AgentDesktop (main export) ────────────────────────────────────────────────

type ViewMode = "board" | "grid";

export function AgentDesktop() {
  const [agents, setAgents] = useState<SquadAgent[]>([]);
  const [view, setView] = useState<ViewMode>("board");

  const fetchSquad = () => {
    fetch("/api/squad")
      .then((res) => res.json())
      .then((data) => { if (Array.isArray(data)) setAgents(data); })
      .catch((err) => console.error("AgentDesktop: failed to fetch squad", err));
  };

  useEffect(() => {
    fetchSquad();
    const es = new EventSource("/api/events");
    es.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (["agent_status_update", "agent_tool", "agent_handoff"].includes(data.type)) {
          fetchSquad();
        }
      } catch { /* ignore */ }
    };
    return () => es.close();
  }, []);

  return (
    <div className="space-y-3">
      {/* view toggle */}
      <div className="flex gap-1 font-mono text-[10px]">
        <button
          onClick={() => setView("board")}
          className={`px-2 py-0.5 rounded border transition-colors ${
            view === "board"
              ? "border-white/30 text-white/80 bg-white/10"
              : "border-white/10 text-white/30 hover:text-white/50"
          }`}
        >
          [board]
        </button>
        <button
          onClick={() => setView("grid")}
          className={`px-2 py-0.5 rounded border transition-colors ${
            view === "grid"
              ? "border-white/30 text-white/80 bg-white/10"
              : "border-white/10 text-white/30 hover:text-white/50"
          }`}
        >
          [grid]
        </button>
      </div>

      {view === "board" ? (
        <TrelloBoard agents={agents} />
      ) : (
        <DesktopGrid agents={agents} />
      )}
    </div>
  );
}
