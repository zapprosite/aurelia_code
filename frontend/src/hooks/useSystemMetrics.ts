import { useState, useEffect } from "react";

export interface SystemMetrics {
  vramUsed: number | null;
  vramTotal: number | null;
  gpuUtil: number | null;
  cpuLoad: number;
  memUsed: number;
  memTotal: number;
  loading: boolean;
  error: string | null;
}

interface RawMetrics {
  vram_used_gib: number | null;
  vram_total_gib: number | null;
  gpu_util_percent: number | null;
  cpu_load: number;
  mem_used_gib: number;
  mem_total_gib: number;
}

const POLL_INTERVAL_MS = 5000;

export function useSystemMetrics(): SystemMetrics {
  const [metrics, setMetrics] = useState<SystemMetrics>({
    vramUsed: null,
    vramTotal: null,
    gpuUtil: null,
    cpuLoad: 0,
    memUsed: 0,
    memTotal: 0,
    loading: true,
    error: null,
  });

  useEffect(() => {
    let cancelled = false;

    const fetchMetrics = () => {
      fetch("/api/metrics")
        .then((res) => {
          if (!res.ok) throw new Error(`HTTP ${res.status}`);
          return res.json() as Promise<RawMetrics>;
        })
        .then((data) => {
          if (!cancelled) {
            setMetrics({
              vramUsed: data.vram_used_gib ?? null,
              vramTotal: data.vram_total_gib ?? null,
              gpuUtil: data.gpu_util_percent ?? null,
              cpuLoad: data.cpu_load ?? 0,
              memUsed: data.mem_used_gib ?? 0,
              memTotal: data.mem_total_gib ?? 0,
              loading: false,
              error: null,
            });
          }
        })
        .catch((err) => {
          if (!cancelled) {
            setMetrics((prev) => ({ ...prev, loading: false, error: String(err) }));
          }
        });
    };

    fetchMetrics();
    const id = setInterval(fetchMetrics, POLL_INTERVAL_MS);
    return () => {
      cancelled = true;
      clearInterval(id);
    };
  }, []);

  return metrics;
}
