import * as React from "react";
import { Brain, Search } from "lucide-react";

interface BrainPoint {
  id: string | number;
  score?: number;
  payload?: Record<string, unknown>;
}

function payloadSummary(payload?: Record<string, unknown>): string {
  if (!payload) return "(sem payload)";
  const vals = Object.values(payload)
    .filter((v) => typeof v === "string")
    .slice(0, 3)
    .map((v) => String(v).slice(0, 120));
  return vals.join(" · ") || JSON.stringify(payload).slice(0, 200);
}

export function BrainSearch() {
  const [query, setQuery] = React.useState("");
  const [results, setResults] = React.useState<BrainPoint[]>([]);
  const [recent, setRecent] = React.useState<BrainPoint[]>([]);
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  // Load recent on mount
  React.useEffect(() => {
    fetch("/api/brain/recent")
      .then((r) => r.json())
      .then((data) => {
        if (Array.isArray(data)) setRecent(data);
      })
      .catch(() => setRecent([]));
  }, []);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;
    setLoading(true);
    setError(null);
    fetch(`/api/brain/search?q=${encodeURIComponent(query.trim())}`)
      .then((r) => r.json())
      .then((data) => {
        if (Array.isArray(data)) setResults(data);
        else setResults([]);
      })
      .catch(() => {
        setError("Erro ao buscar na memória semântica.");
        setResults([]);
      })
      .finally(() => setLoading(false));
  };

  const displayList = results.length > 0 ? results : recent;
  const listLabel = results.length > 0 ? `${results.length} resultado(s) para "${query}"` : "Memórias recentes";

  return (
    <div className="space-y-6">
      {/* Search box */}
      <form onSubmit={handleSearch} className="flex gap-2">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-white/30 pointer-events-none" />
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Buscar memória semântica…"
            className="w-full pl-9 pr-4 py-2 rounded-lg bg-white/5 border border-white/10 text-sm text-white/80 placeholder-white/20 focus:outline-none focus:border-purple-500/50"
          />
        </div>
        <button
          type="submit"
          disabled={loading}
          className="px-4 py-2 rounded-lg bg-purple-600/20 border border-purple-500/30 text-purple-300 text-sm hover:bg-purple-600/30 transition-colors disabled:opacity-40"
        >
          {loading ? "…" : "Buscar"}
        </button>
        {results.length > 0 && (
          <button
            type="button"
            onClick={() => { setResults([]); setQuery(""); }}
            className="px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-white/40 text-sm hover:bg-white/10 transition-colors"
          >
            ✕
          </button>
        )}
      </form>

      {error && (
        <p className="text-red-400 text-sm">{error}</p>
      )}

      {/* Results / Recent list */}
      <div className="space-y-2">
        <p className="text-xs text-white/30 uppercase tracking-widest font-mono">{listLabel}</p>
        {displayList.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-center">
            <Brain className="w-8 h-8 text-white/10 mb-3" />
            <p className="text-sm text-white/20">Nenhuma memória encontrada</p>
          </div>
        ) : (
          displayList.map((point) => (
            <div
              key={String(point.id)}
              className="p-3 rounded-lg bg-white/5 border border-white/10 hover:border-white/20 transition-colors"
            >
              <div className="flex items-start gap-2">
                <span className="text-[10px] font-mono text-white/20 shrink-0 pt-0.5">
                  #{String(point.id).slice(0, 8)}
                </span>
                {point.score !== undefined && (
                  <span className="text-[10px] font-mono text-purple-400/60 shrink-0 pt-0.5">
                    {(point.score * 100).toFixed(0)}%
                  </span>
                )}
                <p className="text-xs text-white/60 break-all leading-relaxed">
                  {payloadSummary(point.payload)}
                </p>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
