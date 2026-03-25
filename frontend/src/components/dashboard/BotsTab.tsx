import * as React from "react";
import { Plus, Bot } from "lucide-react";
import { BotCard, type BotCardData } from "./BotCard";
import { BotDetail } from "./BotDetail";
import { CreateBotModal } from "./CreateBotModal";
import type { FeedItemProps } from "./FeedItem";

interface BotsTabProps {
  feed: FeedItemProps[];
}

export function BotsTab({ feed }: BotsTabProps) {
  const [bots, setBots] = React.useState<BotCardData[]>([]);
  const [selected, setSelected] = React.useState<BotCardData | null>(null);
  const [showCreate, setShowCreate] = React.useState(false);
  const [loading, setLoading] = React.useState(true);

  const fetchBots = () => {
    setLoading(true);
    fetch("/api/bots")
      .then((r) => r.json())
      .then((data: BotCardData[]) => {
        if (Array.isArray(data)) setBots(data);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  };

  React.useEffect(() => {
    fetchBots();
  }, []);

  const handleDelete = (id: string) => {
    if (!confirm(`Remover o bot "${id}" do time?`)) return;
    fetch(`/api/bots/remove?id=${encodeURIComponent(id)}`, { method: "DELETE" })
      .then(() => fetchBots())
      .catch(() => {});
  };

  if (selected) {
    return (
      <BotDetail
        bot={selected}
        feed={feed}
        onBack={() => setSelected(null)}
      />
    );
  }

  return (
    <div className="space-y-6">
      {/* Header row */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2 text-white/40">
          <Bot className="w-4 h-4" />
          <span className="text-xs font-mono uppercase tracking-widest">Team Bots</span>
          <span className="ml-2 text-[10px] text-white/20">{bots.length} configurado{bots.length !== 1 ? "s" : ""}</span>
        </div>
        <button
          onClick={() => setShowCreate(true)}
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-white/10 bg-white/5 text-xs text-white/60 hover:bg-white/10 hover:text-white/90 transition-all"
        >
          <Plus className="w-3.5 h-3.5" />
          Adicionar Bot
        </button>
      </div>

      {/* Grid */}
      {loading ? (
        <div className="text-center text-sm text-white/30 py-12">Carregando bots...</div>
      ) : bots.length === 0 ? (
        <div className="rounded-2xl border border-dashed border-white/10 p-12 text-center">
          <Bot className="w-10 h-10 text-white/20 mx-auto mb-3" />
          <p className="text-sm text-white/40">Nenhum bot configurado.</p>
          <p className="text-xs text-white/25 mt-1">Clique em "Adicionar Bot" para começar.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {bots.map((bot) => (
            <BotCard
              key={bot.id}
              bot={bot}
              onClick={() => setSelected(bot)}
              onDelete={() => handleDelete(bot.id)}
            />
          ))}
        </div>
      )}

      {/* Create modal */}
      {showCreate && (
        <CreateBotModal
          onClose={() => setShowCreate(false)}
          onCreated={fetchBots}
        />
      )}
    </div>
  );
}
