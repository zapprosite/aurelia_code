import { useEffect, useState } from "react";
import * as Icons from "lucide-react";
import { motion } from "framer-motion";
import { Badge } from "../ui/Badge";
import { type SquadAgent } from "./SquadGrid";

// Fixed desk layout: col/row in a 6x4 CSS grid (1-indexed)
const DESK_LAYOUT: Record<string, { col: number; row: number }> = {
  aurelia: { col: 3, row: 1 }, // Central command
  claude:  { col: 2, row: 3 }, // Left station
  codex:   { col: 4, row: 3 }, // Right station
};

// Free desk slots for dynamic agents (col, row pairs not in DESK_LAYOUT)
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
  if (DESK_LAYOUT[key]) {
    return DESK_LAYOUT[key];
  }
  // Assign a free slot to dynamic agents
  for (const slot of FREE_SLOTS) {
    const slotKey = `${slot.col}-${slot.row}`;
    if (!usedSlots.has(slotKey)) {
      usedSlots.add(slotKey);
      return slot;
    }
  }
  // Overflow: place in a default position
  return { col: 6, row: 4 };
}

function statusAnimation(status: string) {
  switch (status) {
    case "busy":
      return {
        opacity: [1, 0.7, 1],
        scale: [1, 1.02, 1],
        transition: { duration: 1.5, repeat: Infinity, ease: "easeInOut" as const },
      };
    case "offline":
      return { opacity: 0.35 };
    default:
      return {};
  }
}

function statusGlow(status: string, color: string) {
  if (status === "online") {
    return `shadow-[0_0_18px_2px] ${color.replace("text-", "shadow-")}/30`;
  }
  return "";
}

export function AgentDesktop() {
  const [agents, setAgents] = useState<SquadAgent[]>([]);

  const fetchSquad = () => {
    fetch("/api/squad")
      .then((res) => res.json())
      .then((data) => {
        if (Array.isArray(data)) setAgents(data);
      })
      .catch((err) => console.error("AgentDesktop: failed to fetch squad", err));
  };

  useEffect(() => {
    fetchSquad();
    const es = new EventSource("/api/events");
    es.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (
          data.type === "agent_status_update" ||
          data.type === "agent_tool" ||
          data.type === "agent_handoff"
        ) {
          fetchSquad();
        }
      } catch {
        // ignore
      }
    };
    return () => es.close();
  }, []);

  // Compute layout: reserve fixed slots, assign free slots to dynamic agents
  const usedSlots = new Set<string>(
    Object.values(DESK_LAYOUT).map((s) => `${s.col}-${s.row}`)
  );
  const agentSlots = agents.map((agent) => ({
    agent,
    slot: agentSlot(agent, usedSlots),
  }));

  return (
    <div className="relative w-full" style={{ minHeight: "520px" }}>
      {/* Office-lines background via CSS gradient */}
      <div
        className="absolute inset-0 pointer-events-none rounded-xl opacity-20"
        style={{
          backgroundImage:
            "repeating-linear-gradient(0deg, transparent, transparent 79px, rgba(255,255,255,0.04) 80px), " +
            "repeating-linear-gradient(90deg, transparent, transparent 79px, rgba(255,255,255,0.04) 80px)",
          backgroundSize: "80px 80px",
        }}
      />

      {/* CSS Grid 6 columns x 4 rows */}
      <div
        className="relative grid gap-4 p-4"
        style={{
          gridTemplateColumns: "repeat(6, 1fr)",
          gridTemplateRows: "repeat(4, auto)",
        }}
      >
        {agentSlots.map(({ agent, slot }, i) => {
          const Icon = getIcon(agent.icon);
          const anim = statusAnimation(agent.status);
          const glowClass = statusGlow(agent.status, agent.color);

          return (
            <motion.div
              key={agent.id}
              initial={{ opacity: 0, scale: 0.85 }}
              animate={{ opacity: 1, scale: 1, ...anim }}
              transition={{ delay: i * 0.08 }}
              style={{
                gridColumn: slot.col,
                gridRow: slot.row,
              }}
            >
              <div
                className={`relative rounded-xl border border-white/10 bg-white/5 p-4 cursor-default hover:border-white/25 transition-all overflow-hidden ${glowClass}`}
              >
                {/* Color glow blob */}
                <div
                  className={`absolute top-0 right-0 w-20 h-20 blur-3xl opacity-10 transition-opacity hover:opacity-20 ${agent.color.replace(
                    "text-",
                    "bg-"
                  )}`}
                />

                <div className="flex items-center justify-between mb-3 relative z-10">
                  <div
                    className={`p-2 rounded-lg bg-white/5 border border-white/5 ${agent.color}`}
                  >
                    <Icon className="w-4 h-4" />
                  </div>
                  <Badge
                    variant={
                      agent.status === "online"
                        ? "success"
                        : agent.status === "busy"
                        ? "default"
                        : "outline"
                    }
                    className="uppercase text-[9px]"
                  >
                    {agent.status}
                  </Badge>
                </div>

                <div className="relative z-10">
                  <h4 className="font-bold text-white/90 tracking-tight text-sm">
                    {agent.name}
                  </h4>
                  <p className="text-[10px] text-white/40 font-medium mb-2">
                    {agent.role}
                  </p>
                  <div className="space-y-1">
                    <div className="flex justify-between text-[9px] text-white/30 font-mono">
                      <span>Load</span>
                      <span>{agent.load}%</span>
                    </div>
                    <div className="h-0.5 w-full bg-white/5 rounded-full overflow-hidden">
                      <motion.div
                        initial={{ width: 0 }}
                        animate={{ width: `${agent.load}%` }}
                        className={`h-full ${agent.color.replace("text-", "bg-")}`}
                      />
                    </div>
                  </div>
                </div>
              </div>
            </motion.div>
          );
        })}
      </div>
    </div>
  );
}
