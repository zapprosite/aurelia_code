import * as React from "react";
import { ArrowLeft, Crown, Thermometer, ClipboardCheck, Calendar, Bot } from "lucide-react";
import { FeedItem, type FeedItemProps } from "./FeedItem";
import type { BotCardData } from "./BotCard";

interface BotDetailProps {
  bot: BotCardData;
  feed: FeedItemProps[];
  onBack: () => void;
}

const PERSONA_META: Record<string, { icon: React.ReactNode; color: string }> = {
  "aurelia-leader":   { icon: <Crown className="w-5 h-5" />,          color: "text-purple-400" },
  "hvac-sales":       { icon: <Thermometer className="w-5 h-5" />,    color: "text-blue-400" },
  "project-manager":  { icon: <ClipboardCheck className="w-5 h-5" />, color: "text-yellow-400" },
  "life-organizer":   { icon: <Calendar className="w-5 h-5" />,        color: "text-green-400" },
};

export function BotDetail({ bot, feed, onBack }: BotDetailProps) {
  const meta = PERSONA_META[bot.persona_id] ?? { icon: <Bot className="w-5 h-5" />, color: "text-white/60" };

  // Filter feed by bot_id (the feed items may carry a bot_id from SSE event)
  const botFeed = feed.filter((f) => {
    // FeedItem doesn't have bot_id by default; we enrich via agent name heuristic
    return (f as any).botId === bot.id || f.agent?.toLowerCase() === bot.id.toLowerCase();
  });

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <button
          className="text-white/40 hover:text-white/80 transition-colors"
          onClick={onBack}
        >
          <ArrowLeft className="w-5 h-5" />
        </button>
        <div className={`${meta.color}`}>{meta.icon}</div>
        <div>
          <h2 className="text-base font-semibold text-white/90">{bot.name}</h2>
          <span className="text-xs font-mono text-white/40">@{bot.id}</span>
        </div>
        {/* Status */}
        <div className="ml-auto flex items-center gap-2">
          <span className="relative flex h-2.5 w-2.5">
            {bot.running ? (
              <>
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75" />
                <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500" />
              </>
            ) : (
              <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-white/20" />
            )}
          </span>
          <span className="text-xs text-white/40">{bot.running ? "Online" : "Offline"}</span>
        </div>
      </div>

      {/* Info grid */}
      <div className="grid grid-cols-2 gap-4">
        <div className="rounded-xl border border-white/10 bg-white/5 p-4">
          <div className="text-[10px] font-mono uppercase tracking-widest text-white/30 mb-1">Persona</div>
          <div className="text-sm text-white/80">{bot.persona_id || "—"}</div>
        </div>
        <div className="rounded-xl border border-white/10 bg-white/5 p-4">
          <div className="text-[10px] font-mono uppercase tracking-widest text-white/30 mb-1">Área de Foco</div>
          <div className="text-sm text-white/80">{bot.focus_area || "—"}</div>
        </div>
      </div>

      {/* Activity feed */}
      <div>
        <div className="text-[10px] font-mono uppercase tracking-widest text-white/30 mb-3">
          Atividade Recente
        </div>
        {botFeed.length === 0 ? (
          <div className="rounded-xl border border-white/10 bg-white/5 p-6 text-center text-sm text-white/30">
            Nenhuma atividade registrada para este bot ainda.
          </div>
        ) : (
          <div className="space-y-3">
            {botFeed.slice(0, 20).map((item) => (
              <FeedItem key={item.id} {...item} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
