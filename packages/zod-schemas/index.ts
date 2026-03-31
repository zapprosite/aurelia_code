import { z } from 'zod';

/**
 * 🛰️ Sentinel Core Schemas — SOTA 2026
 * Autoridade: Aurélia Governance
 */

export const SentinelEventSchema = z.object({
  id: z.string().uuid(),
  timestamp: z.string().datetime(),
  type: z.enum(['SECURITY_ALARM', 'DB_HEALTH_CHECK', 'RATE_LIMIT_HIT', 'SECRET_LEAK_PREVENTED']),
  severity: z.enum(['LOW', 'MEDIUM', 'HIGH', 'CRITICAL']),
  metadata: z.record(z.any()),
});

export type SentinelEvent = z.infer<typeof SentinelEventSchema>;

export const DatabaseAuditSchema = z.object({
  status: z.enum(['OK', 'WARNING', 'ERROR']),
  last_check: z.string().datetime(),
  slow_queries_count: z.number(),
  migrations_synced: z.boolean(),
});
