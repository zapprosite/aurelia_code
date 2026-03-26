import * as React from "react";
import { Activity, Bot, Boxes, HardDrive, RefreshCcw, Server, ShieldCheck, TerminalSquare } from "lucide-react";
import { Badge } from "../ui/Badge";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "../ui/Card";

type HomelabSnapshot = {
  containers: Array<{ name: string; status: string; ports: string; up: boolean }>;
  health: Array<{ name: string; url: string; ok: boolean; code: number }>;
  snapshots: Array<{ name: string; creation: string }>;
  changelog: string;
  agent_state: { lastAction: string; timestamp: string; status: string; agent: string; model: string };
  zpool_status: string;
  disk_usage: {
    raw: string;
    filesystem: string;
    size: string;
    used: string;
    available: string;
    used_percent: number;
    mount: string;
  };
  timestamp: string;
};

function formatTimestamp(value?: string) {
  if (!value) return "agora";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat("pt-BR", {
    dateStyle: "short",
    timeStyle: "medium"
  }).format(date);
}

function StatusDot({ ok }: { ok: boolean }) {
  return (
    <span
      className={[
        "inline-flex h-2.5 w-2.5 rounded-full",
        ok ? "bg-emerald-400 shadow-[0_0_12px_rgba(74,222,128,0.65)]" : "bg-rose-400 shadow-[0_0_12px_rgba(251,113,133,0.55)]"
      ].join(" ")}
    />
  );
}

function SummaryCard({
  icon,
  label,
  value,
  detail,
  tone = "neutral"
}: {
  icon: React.ReactNode;
  label: string;
  value: string;
  detail: string;
  tone?: "neutral" | "good" | "warn";
}) {
  const toneClass =
    tone === "good"
      ? "border-emerald-500/20 bg-emerald-500/8"
      : tone === "warn"
        ? "border-amber-500/20 bg-amber-500/8"
        : "border-white/10 bg-white/5";

  return (
    <Card className={toneClass}>
      <CardHeader className="pb-3">
        <div className="flex items-center gap-2 text-white/45">
          {icon}
          <CardDescription className="text-[11px] font-mono uppercase tracking-[0.2em] text-white/45">
            {label}
          </CardDescription>
        </div>
      </CardHeader>
      <CardContent className="space-y-1">
        <div className="text-3xl font-semibold tracking-tight text-white/92">{value}</div>
        <p className="text-sm text-white/45">{detail}</p>
      </CardContent>
    </Card>
  );
}

function DataTable({ headers, rows }: { headers: string[]; rows: React.ReactNode }) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full min-w-[640px] text-sm">
        <thead>
          <tr className="border-b border-white/10 text-left text-[11px] uppercase tracking-[0.18em] text-white/35">
            {headers.map((header) => (
              <th key={header} className="px-4 py-3 font-medium">{header}</th>
            ))}
          </tr>
        </thead>
        <tbody>{rows}</tbody>
      </table>
    </div>
  );
}

export function HomelabTab() {
  const [snapshot, setSnapshot] = React.useState<HomelabSnapshot | null>(null);
  const [loading, setLoading] = React.useState(true);
  const [refreshing, setRefreshing] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  const loadSnapshot = React.useEffectEvent(async (background = false) => {
    if (!background) {
      setLoading(snapshot === null);
    }
    setRefreshing(true);

    try {
      const response = await fetch("/api/homelab", { cache: "no-store" });
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
      }
      const data = await response.json() as HomelabSnapshot;
      React.startTransition(() => {
        setSnapshot(data);
        setError(null);
      });
    } catch (err) {
      const message = err instanceof Error ? err.message : "Falha ao carregar monitor";
      setError(message);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  });

  React.useEffect(() => {
    void loadSnapshot(false);
    const timer = window.setInterval(() => {
      void loadSnapshot(true);
    }, 30000);
    return () => window.clearInterval(timer);
  }, []);

  const containers = snapshot?.containers ?? [];
  const health = snapshot?.health ?? [];
  const snapshots = snapshot?.snapshots ?? [];
  const containersUp = containers.filter((item) => item.up).length;
  const containersDown = containers.length - containersUp;
  const healthUp = health.filter((item) => item.ok).length;
  const healthDown = health.length - healthUp;
  const diskPercent = snapshot?.disk_usage.used_percent ?? 0;

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center gap-3 text-white/40">
        <div className="flex items-center gap-2">
          <Server className="h-4 w-4" />
          <span className="text-xs font-mono uppercase tracking-[0.28em]">Homelab Control Deck</span>
        </div>
        <Badge variant="outline" className="text-[10px] uppercase tracking-[0.2em] text-white/45">
          auto 30s
        </Badge>
        <span className="text-xs text-white/30">
          atualizado em {formatTimestamp(snapshot?.timestamp)}
        </span>
        {error && (
          <span className="text-xs text-amber-300/80">última falha: {error}</span>
        )}
        <button
          onClick={() => {
            void loadSnapshot(false);
          }}
          disabled={refreshing}
          className="ml-auto inline-flex items-center gap-2 rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-xs font-medium text-white/70 transition-colors hover:border-white/30 hover:text-white disabled:cursor-not-allowed disabled:opacity-50"
          title="Atualizar dados do homelab"
        >
          <RefreshCcw className={["h-3.5 w-3.5", refreshing ? "animate-spin" : ""].join(" ")} />
          atualizar
        </button>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
        <SummaryCard
          icon={<Boxes className="h-4 w-4" />}
          label="Containers"
          value={`${containersUp}/${containers.length || 0}`}
          detail={containersDown > 0 ? `${containersDown} com falha ou desligados` : "todos operacionais"}
          tone={containersDown > 0 ? "warn" : "good"}
        />
        <SummaryCard
          icon={<ShieldCheck className="h-4 w-4" />}
          label="Endpoints"
          value={`${healthUp}/${health.length || 0}`}
          detail={healthDown > 0 ? `${healthDown} indisponível` : "health checks OK"}
          tone={healthDown > 0 ? "warn" : "good"}
        />
        <SummaryCard
          icon={<Activity className="h-4 w-4" />}
          label="Snapshots ZFS"
          value={`${snapshots.length}`}
          detail="cauda mais recente do histórico"
        />
        <SummaryCard
          icon={<HardDrive className="h-4 w-4" />}
          label="Disco /srv"
          value={snapshot?.disk_usage.used || "N/A"}
          detail={snapshot?.disk_usage.raw || "sem leitura"}
          tone={diskPercent >= 85 ? "warn" : "neutral"}
        />
      </div>

      {loading && !snapshot && (
        <Card>
          <CardContent className="p-10 text-center text-sm text-white/45">
            carregando o estado operacional do homelab...
          </CardContent>
        </Card>
      )}

      {!loading && !snapshot && (
        <Card className="border-amber-500/20 bg-amber-500/8">
          <CardContent className="p-10 text-center text-sm text-amber-100/85">
            não consegui montar o painel do homelab agora.
          </CardContent>
        </Card>
      )}

      {snapshot && (
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
          <Card className="lg:col-span-2">
            <CardHeader className="pb-4">
              <CardTitle className="flex items-center gap-2 text-base">
                <Boxes className="h-4 w-4 text-primary" />
                Containers
              </CardTitle>
              <CardDescription>estado bruto do Docker local, sem proxy externo</CardDescription>
            </CardHeader>
            <CardContent className="pt-0">
              <DataTable
                headers={["Nome", "Estado", "Status", "Portas"]}
                rows={containers.map((item) => (
                  <tr key={item.name} className="border-b border-white/6 text-white/78 last:border-b-0">
                    <td className="px-4 py-3 font-medium text-white/88">{item.name}</td>
                    <td className="px-4 py-3">
                      <div className="inline-flex items-center gap-2">
                        <StatusDot ok={item.up} />
                        <span className={item.up ? "text-emerald-300/90" : "text-rose-300/90"}>
                          {item.up ? "up" : "off"}
                        </span>
                      </div>
                    </td>
                    <td className="px-4 py-3 text-white/58">{item.status}</td>
                    <td className="px-4 py-3 font-mono text-xs text-white/45">{item.ports || "—"}</td>
                  </tr>
                ))}
              />
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-4">
              <CardTitle className="flex items-center gap-2 text-base">
                <ShieldCheck className="h-4 w-4 text-primary" />
                Health + Agente
              </CardTitle>
              <CardDescription>probes locais e estado do agente que abastecia o VRV</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6 pt-0">
              <div className="space-y-3">
                {health.map((item) => (
                  <div key={item.name} className="flex items-start justify-between gap-3 rounded-lg border border-white/8 bg-white/[0.03] px-3 py-3">
                    <div>
                      <div className="flex items-center gap-2 text-sm font-medium text-white/85">
                        <StatusDot ok={item.ok} />
                        {item.name}
                      </div>
                      <div className="mt-1 text-xs text-white/35">{item.url}</div>
                    </div>
                    <div className="text-right text-xs text-white/45">
                      HTTP {item.code || 0}
                    </div>
                  </div>
                ))}
              </div>

              <div className="rounded-xl border border-white/8 bg-white/[0.03] p-4">
                <div className="mb-3 flex items-center gap-2 text-sm font-medium text-white/85">
                  <Bot className="h-4 w-4 text-primary" />
                  Agente
                </div>
                <dl className="grid grid-cols-2 gap-x-3 gap-y-2 text-sm">
                  <dt className="text-white/35">nome</dt>
                  <dd className="text-white/82">{snapshot.agent_state.agent}</dd>
                  <dt className="text-white/35">modelo</dt>
                  <dd className="text-white/82">{snapshot.agent_state.model}</dd>
                  <dt className="text-white/35">ação</dt>
                  <dd className="text-white/82">{snapshot.agent_state.lastAction}</dd>
                  <dt className="text-white/35">status</dt>
                  <dd className="text-white/82">{snapshot.agent_state.status}</dd>
                  <dt className="text-white/35">timestamp</dt>
                  <dd className="text-white/82">{formatTimestamp(snapshot.agent_state.timestamp)}</dd>
                </dl>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-4">
              <CardTitle className="flex items-center gap-2 text-base">
                <HardDrive className="h-4 w-4 text-primary" />
                Disco /srv
              </CardTitle>
              <CardDescription>leitura direta do host onde a Aurelia roda</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4 pt-0">
              <div className="rounded-xl border border-white/8 bg-white/[0.03] p-4">
                <div className="mb-2 text-2xl font-semibold text-white/90">
                  {snapshot.disk_usage.used || "N/A"}
                </div>
                <div className="mb-3 text-sm text-white/45">{snapshot.disk_usage.raw || "sem leitura"}</div>
                <div className="h-2 overflow-hidden rounded-full bg-white/8">
                  <div
                    className={[
                      "h-full rounded-full",
                      diskPercent >= 85 ? "bg-amber-400" : "bg-primary"
                    ].join(" ")}
                    style={{ width: `${Math.min(Math.max(diskPercent, 0), 100)}%` }}
                  />
                </div>
                <dl className="mt-4 grid grid-cols-2 gap-x-3 gap-y-2 text-sm">
                  <dt className="text-white/35">filesystem</dt>
                  <dd className="font-mono text-xs text-white/72">{snapshot.disk_usage.filesystem || "-"}</dd>
                  <dt className="text-white/35">disponível</dt>
                  <dd className="text-white/82">{snapshot.disk_usage.available || "-"}</dd>
                  <dt className="text-white/35">tamanho</dt>
                  <dd className="text-white/82">{snapshot.disk_usage.size || "-"}</dd>
                  <dt className="text-white/35">mount</dt>
                  <dd className="text-white/82">{snapshot.disk_usage.mount || "-"}</dd>
                </dl>
              </div>

              <div className="rounded-xl border border-white/8 bg-white/[0.03] p-4">
                <div className="mb-3 flex items-center gap-2 text-sm font-medium text-white/85">
                  <TerminalSquare className="h-4 w-4 text-primary" />
                  ZPool
                </div>
                <pre className="overflow-x-auto whitespace-pre-wrap font-mono text-xs leading-6 text-white/62">
                  {snapshot.zpool_status || "ZFS pool não disponível"}
                </pre>
              </div>
            </CardContent>
          </Card>

          <Card className="lg:col-span-2">
            <CardHeader className="pb-4">
              <CardTitle className="flex items-center gap-2 text-base">
                <Activity className="h-4 w-4 text-primary" />
                Snapshots ZFS
              </CardTitle>
              <CardDescription>cauda dos snapshots mais recentes do pool</CardDescription>
            </CardHeader>
            <CardContent className="pt-0">
              <DataTable
                headers={["Snapshot", "Criação"]}
                rows={snapshots.map((item) => (
                  <tr key={`${item.name}-${item.creation}`} className="border-b border-white/6 text-white/78 last:border-b-0">
                    <td className="px-4 py-3 font-mono text-xs text-white/72">{item.name}</td>
                    <td className="px-4 py-3 text-white/52">{item.creation}</td>
                  </tr>
                ))}
              />
            </CardContent>
          </Card>

          <Card className="lg:col-span-3">
            <CardHeader className="pb-4">
              <CardTitle className="flex items-center gap-2 text-base">
                <TerminalSquare className="h-4 w-4 text-primary" />
                Changelog
              </CardTitle>
              <CardDescription>últimas entradas de `/srv/ops/ai-governance/logs/CHANGE_LOG.txt`</CardDescription>
            </CardHeader>
            <CardContent className="pt-0">
              <pre className="max-h-[32rem] overflow-auto rounded-xl border border-white/8 bg-[#0a0f18] p-4 font-mono text-xs leading-6 text-cyan-100/75">
                {snapshot.changelog || "Nenhum changelog encontrado"}
              </pre>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  );
}
