import * as React from "react";
import { Crown, Thermometer, ClipboardCheck, Calendar, Bot, Trash2 } from "lucide-react";

export interface BotCardData {
  id: string;
  name: string;
  persona_id: string;
  focus_area: string;
  enabled: boolean;
  running: boolean;
}

interface BotCardProps {
  bot: BotCardData;
  onClick: () => void;
  onDelete: () => void;
}

const PERSONA_META: Record<string, { icon: React.ReactNode; color: string; label: string }> = {
  "aurelia-leader":   { icon: <Crown className="w-5 h-5" />,         color: "text-purple-400",  label: "Líder" },
  "hvac-sales":       { icon: <Thermometer className="w-5 h-5" />,   color: "text-blue-400",    label: "Vendas HVAC-R" },
  "project-manager":  { icon: <ClipboardCheck className="w-5 h-5" />, color: "text-yellow-400", label: "Gestor de Obras" },
  "life-organizer":   { icon: <Calendar className="w-5 h-5" />,       color: "text-green-400",  label: "Vida & Agenda" },
};

export function BotCard({ bot, onClick, onDelete }: BotCardProps) {
  const meta = PERSONA_META[bot.persona_id] ?? {
    icon: <Bot className="w-5 h-5" />,
    color: "text-white/60",
    label: bot.persona_id || "Bot",
  };

  return (
    <div
      className="relative group cursor-pointer rounded-xl border border-white/10 bg-white/5 backdrop-blur p-5 hover:border-white/20 hover:bg-white/8 transition-all duration-200"
      onClick={onClick}
    >
      {/* Status dot */}
      <span className="absolute top-4 right-4 flex h-2.5 w-2.5">
        {bot.running ? (
          <>
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75" />
            <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500" />
          </>
        ) : (
          <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-white/20" />
        )}
      </span>

      {/* Delete button */}
      <button
        className="absolute bottom-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity text-white/30 hover:text-red-400"
        onClick={(e) => {
          e.stopPropagation();
          onDelete();
        }}
        title="Remover bot"
      >
        <Trash2 className="w-4 h-4" />
      </button>

      {/* Icon + name */}
      <div className="flex items-center gap-3 mb-3">
        <div className={`${meta.color}`}>{meta.icon}</div>
        <div>
          <div className="text-sm font-semibold text-white/90">{bot.name}</div>
          <div className={`text-[11px] font-mono ${meta.color} opacity-80`}>{meta.label}</div>
        </div>
      </div>

      {/* Focus area */}
      {bot.focus_area && (
        <p className="text-[11px] text-white/40 leading-relaxed line-clamp-2">{bot.focus_area}</p>
      )}

      {/* ID chip */}
      <div className="mt-3 inline-flex items-center px-2 py-0.5 rounded bg-white/5 border border-white/10">
        <span className="text-[10px] font-mono text-white/30">@{bot.id}</span>
      </div>
    </div>
  );
}
