import * as React from "react";
import { Activity, Bot, Boxes, HardDrive, RefreshCcw, Server, ShieldCheck, TerminalSquare } from "lucide-react";
import { Badge } from "../ui/Badge";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "../ui/Card";

type HomelabServiceStatus = "healthy" | "degraded" | "offline";

type HomelabServiceContainer = {
  name?: string;
  status?: string;
  ports?: string;
  up?: boolean;
};

type HomelabService = {
  id?: string;
  key?: string;
  label?: string;
  name?: string;
  slug?: string;
  service?: string;
  status?: string;
  state?: string;
  health?: string;
  summary?: string;
  detail?: string;
  description?: string;
  message?: string;
  reason?: string;
  counts?: {
    total?: number;
    up?: number;
    restarting?: number;
    exited?: number;
    dead?: number;
    other?: number;
  };
  url?: string;
  endpoint?: string;
  container?: string;
  containers?: Array<string | HomelabServiceContainer> | string;
  ok?: boolean;
  up?: boolean;
  degraded?: boolean;
};

type HomelabSnapshot = {
  services?: HomelabService[] | Record<string, HomelabService>;
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

type ServiceDefinition = {
  key: "supabase" | "qdrant" | "caprover";
  label: string;
  icon: React.ReactNode;
  fallbackSummary: string;
  fallbackDetails: string;
};

type ServiceView = ServiceDefinition & {
  status: HomelabServiceStatus;
  summary: string;
  detail: string;
  source: string;
};

const CANONICAL_SERVICES: ServiceDefinition[] = [
  {
    key: "supabase",
    label: "Supabase",
    icon: <ShieldCheck className="h-4 w-4" />,
    fallbackSummary: "registro canônico e camada de dados local.",
    fallbackDetails: "infra de dados e auth do stack."
  },
  {
    key: "qdrant",
    label: "Qdrant",
    icon: <Activity className="h-4 w-4" />,
    fallbackSummary: "índice vetorial e busca semântica.",
    fallbackDetails: "vetores e recuperação semântica."
  },
  {
    key: "caprover",
    label: "CapRover",
    icon: <Server className="h-4 w-4" />,
    fallbackSummary: "plataforma de deploy e exposição dos apps.",
    fallbackDetails: "orquestração e publicação de serviços."
  }
];

function formatTimestamp(value?: string) {
  if (!value) return "agora";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat("pt-BR", {
    dateStyle: "short",
    timeStyle: "medium"
  }).format(date);
}

function toList<T>(value: T[] | Record<string, T> | undefined): T[] {
  if (!value) return [];
  return Array.isArray(value) ? value : Object.values(value);
}

function statusLabel(status: HomelabServiceStatus) {
  switch (status) {
    case "healthy":
      return "healthy";
    case "degraded":
      return "degraded";
    case "offline":
    default:
      return "offline";
  }
}

function statusStyles(status: HomelabServiceStatus) {
  switch (status) {
    case "healthy":
      return {
        card: "border-emerald-500/20 bg-emerald-500/8",
        dot: "bg-emerald-400 shadow-[0_0_12px_rgba(74,222,128,0.65)]",
        badge: "border-emerald-500/20 bg-emerald-500/12 text-emerald-200",
        text: "text-emerald-300/90"
      };
    case "degraded":
      return {
        card: "border-amber-500/20 bg-amber-500/8",
        dot: "bg-amber-400 shadow-[0_0_12px_rgba(251,191,36,0.55)]",
        badge: "border-amber-500/20 bg-amber-500/12 text-amber-100",
        text: "text-amber-300/90"
      };
    case "offline":
    default:
      return {
        card: "border-rose-500/20 bg-rose-500/8",
        dot: "bg-rose-400 shadow-[0_0_12px_rgba(251,113,133,0.55)]",
        badge: "border-rose-500/20 bg-rose-500/12 text-rose-100",
        text: "text-rose-300/90"
      };
  }
}

function normalizeStatus(raw: unknown, fallbackUp?: boolean): HomelabServiceStatus {
  if (typeof raw === "boolean") {
    return raw ? "healthy" : "offline";
  }

  if (typeof raw === "string") {
    const value = raw.trim().toLowerCase();
    if (!value) {
      return fallbackUp === undefined ? "offline" : fallbackUp ? "healthy" : "offline";
    }
    if (["healthy", "ok", "up", "online", "ready", "pass", "passed"].includes(value) || value.startsWith("up") || value.includes("healthy") || value.includes("ready")) {
      return "healthy";
    }
    if (["degraded", "warning", "warn", "partial"].includes(value) || value.includes("degrad") || value.includes("warn") || value.includes("partial")) {
      return "degraded";
    }
    if (["offline", "down", "error", "fail", "failed", "unavailable", "missing", "stopped"].includes(value) || value.includes("offline") || value.includes("down") || value.includes("error") || value.includes("fail")) {
      return "offline";
    }
  }

  if (fallbackUp !== undefined) {
    return fallbackUp ? "healthy" : "offline";
  }

  return "offline";
}

function findText(...values: Array<unknown>) {
  for (const value of values) {
    if (typeof value === "string" && value.trim()) return value.trim();
  }
  return "";
}

function normalizeContainers(value: HomelabService["containers"]) {
  if (!value) return [] as string[];
  if (typeof value === "string") {
    return value
      .split(/[;,]/)
      .map((item) => item.trim())
      .filter(Boolean);
  }
  return value
    .map((item) => {
      if (typeof item === "string") return item.trim();
      return item.name?.trim() ?? "";
    })
    .filter(Boolean);
}

function isSupabaseLike(value: string) {
  const normalized = value.toLowerCase();
  return normalized.includes("supabase") || normalized.includes("postgres");
}

function isQdrantLike(value: string) {
  const normalized = value.toLowerCase();
  return normalized.includes("qdrant");
}

function isCaproverLike(value: string) {
  const normalized = value.toLowerCase();
  return normalized.includes("caprover") || normalized.includes("captain");
}

function matchesService(service: HomelabService, key: ServiceDefinition["key"]) {
  const haystack = [service.id, service.key, service.label, service.name, service.slug, service.service, service.container, service.url, service.endpoint]
    .filter(Boolean)
    .join(" ")
    .toLowerCase();

  switch (key) {
    case "supabase":
      return isSupabaseLike(haystack);
    case "qdrant":
      return isQdrantLike(haystack);
    case "caprover":
      return isCaproverLike(haystack);
  }
}

function serviceFallback(status: HomelabServiceStatus, label: string, fallbackSummary: string) {
  if (status === "healthy") return fallbackSummary;
  if (status === "degraded") return `${label} está vivo, mas pede atenção.`;
  return `${label} está offline ou sem resposta do snapshot.`;
}

function buildServiceView(
  definition: ServiceDefinition,
  services: HomelabService[],
  health: Array<{ name: string; url: string; ok: boolean; code: number }>,
  containers: Array<{ name: string; status: string; ports: string; up: boolean }>
): ServiceView {
  const rawService = services.find((service) => matchesService(service, definition.key));
  const healthMatch = health.find((item) => {
    const name = item.name.toLowerCase();
    switch (definition.key) {
      case "supabase":
        return name.includes("supabase");
      case "qdrant":
        return name.includes("qdrant");
      case "caprover":
        return name.includes("caprover") || name.includes("captain");
    }
  });

  const rawContainerNames = normalizeContainers(rawService?.containers);
  const fallbackContainer =
    rawService?.container
    ?? rawContainerNames[0]
    ?? containers.find((item) => matchesService({ name: item.name, container: item.name, status: item.status }, definition.key))?.name
    ?? "";
  const explicitStatus = rawService?.degraded ? "degraded" : (rawService?.status ?? rawService?.state ?? rawService?.health ?? rawService?.ok ?? rawService?.up);
  const status = normalizeStatus(explicitStatus, rawContainerNames.length > 0 ? rawContainerNames.some((name) => containers.some((container) => container.name === name && container.up)) : undefined);
  const source = findText(rawService?.url, rawService?.endpoint, healthMatch?.url, fallbackContainer);
  const summary = findText(rawService?.summary, rawService?.message, rawService?.detail, rawService?.description, rawService?.reason)
    || serviceFallback(status, definition.label, definition.fallbackSummary);
  const detail = findText(rawService?.detail, rawService?.description, rawService?.reason)
    || (rawService?.counts?.total ? `${rawService.counts.up ?? 0}/${rawService.counts.total} up${rawService.counts.restarting ? ` · ${rawService.counts.restarting} restarting` : ""}${rawService.counts.exited ? ` · ${rawService.counts.exited} exited` : ""}${rawService.counts.dead ? ` · ${rawService.counts.dead} dead` : ""}` : "")
    || (rawContainerNames.length > 0 ? `containers: ${rawContainerNames.slice(0, 2).join(", ")}${rawContainerNames.length > 2 ? ` +${rawContainerNames.length - 2}` : ""}` : "")
    || (source ? `fonte: ${source}` : definition.fallbackDetails);

  return {
    ...definition,
    status,
    summary,
    detail,
    source: source || definition.fallbackDetails
  };
}

function ServiceStatusDot({ status }: { status: HomelabServiceStatus }) {
  return <span className={["inline-flex h-2.5 w-2.5 rounded-full", statusStyles(status).dot].join(" ")} />;
}

function ServiceCard({ service }: { service: ServiceView }) {
  const styles = statusStyles(service.status);
  return (
    <Card className={["overflow-hidden border-l-2", styles.card].join(" ")}>
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between gap-3">
          <div className="flex items-center gap-2 text-white/45">
            {service.icon}
            <CardTitle className="text-base text-white/88">{service.label}</CardTitle>
          </div>
          <Badge variant="outline" className={["flex items-center gap-2 text-[10px] uppercase tracking-[0.2em]", styles.badge].join(" ")}>
            <ServiceStatusDot status={service.status} />
            <span>{statusLabel(service.status)}</span>
          </Badge>
        </div>
        <CardDescription className="mt-2 text-[11px] font-mono uppercase tracking-[0.2em] text-white/35">
          fonte principal: services
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-3">
        <p className="text-sm leading-6 text-white/80">{service.summary}</p>
        <div className="space-y-2 rounded-xl border border-white/8 bg-white/[0.03] px-3 py-3 text-xs text-white/50">
          <div className="font-medium text-white/72">{service.detail}</div>
          <div className="flex flex-wrap gap-2">
            <span className="rounded-full border border-white/10 bg-black/20 px-2 py-1 font-mono text-[10px] uppercase tracking-[0.18em] text-white/45">
              {service.source}
            </span>
            <span className="rounded-full border border-white/10 bg-black/20 px-2 py-1 font-mono text-[10px] uppercase tracking-[0.18em] text-white/45">
              {statusLabel(service.status)}
            </span>
          </div>
        </div>
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

  const services = toList(snapshot?.services);
  const containers = snapshot?.containers ?? [];
  const health = snapshot?.health ?? [];
  const snapshots = snapshot?.snapshots ?? [];
  const diskPercent = snapshot?.disk_usage.used_percent ?? 0;
  const canonicalServices = CANONICAL_SERVICES.map((definition) =>
    buildServiceView(definition, services, health, containers)
  );
  const serviceStatuses = canonicalServices.reduce<Record<HomelabServiceStatus, number>>((acc, service) => {
    acc[service.status] += 1;
    return acc;
  }, { healthy: 0, degraded: 0, offline: 0 });
  const hasDegradation = serviceStatuses.degraded > 0 || serviceStatuses.offline > 0;

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

      <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
        {canonicalServices.map((service) => (
          <ServiceCard key={service.key} service={service} />
        ))}
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
                Containers brutos
              </CardTitle>
              <CardDescription>detalhe auxiliar do Docker local; a leitura principal vem de `services`</CardDescription>
            </CardHeader>
            <CardContent className="pt-0">
              <DataTable
                headers={["Nome", "Estado", "Status", "Portas"]}
                rows={containers.map((item) => (
                  <tr key={item.name} className="border-b border-white/6 text-white/78 last:border-b-0">
                    <td className="px-4 py-3 font-medium text-white/88">{item.name}</td>
                    <td className="px-4 py-3">
                      <div className="inline-flex items-center gap-2">
                        <ServiceStatusDot status={item.up ? "healthy" : "offline"} />
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
                Checks auxiliares + Agente
              </CardTitle>
              <CardDescription>probes locais e estado do agente; suporte, não fonte principal</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6 pt-0">
              <div className="space-y-2 rounded-xl border border-white/8 bg-white/[0.03] p-3">
                <div className="flex items-center justify-between text-[11px] uppercase tracking-[0.18em] text-white/35">
                  <span>serviços canônicos</span>
                  <span>{serviceStatuses.healthy}/3 healthy</span>
                </div>
                <div className="flex h-2 overflow-hidden rounded-full bg-white/8">
                  <div className="bg-emerald-400" style={{ width: `${(serviceStatuses.healthy / 3) * 100}%` }} />
                  <div className="bg-amber-400" style={{ width: `${(serviceStatuses.degraded / 3) * 100}%` }} />
                  <div className="bg-rose-400" style={{ width: `${(serviceStatuses.offline / 3) * 100}%` }} />
                </div>
                <div className="text-xs text-white/45">
                  {hasDegradation ? `${serviceStatuses.degraded} degraded / ${serviceStatuses.offline} offline` : "todos os serviços canônicos estão healthy"}
                </div>
              </div>

              <div className="space-y-3">
                {health.map((item) => (
                  <div key={item.name} className="flex items-start justify-between gap-3 rounded-lg border border-white/8 bg-white/[0.03] px-3 py-3">
                    <div>
                      <div className="flex items-center gap-2 text-sm font-medium text-white/85">
                        <ServiceStatusDot status={item.ok ? "healthy" : "offline"} />
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
