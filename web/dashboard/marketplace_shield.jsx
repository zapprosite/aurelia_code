import React, { useMemo } from 'react';
import { FramerMotion } from 'framer-motion'; // Assuming available in 2026 env

const MarketplaceView = ({ reputation }) => {
  const sortedAgents = useMemo(() => 
    Object.entries(reputation).sort((a, b) => b[1] - a[1]), 
    [reputation]
  );

  return (
    <div className="marketplace-container glass-morphism p-6 rounded-2xl shadow-[0_20px_50px_rgba(0,0,0,0.5)] border border-white/5 bg-slate-900/40 backdrop-blur-xl">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-2xl font-black bg-clip-text text-transparent bg-gradient-to-r from-emerald-400 to-cyan-400 tracking-tight">
          Talent Agora
        </h3>
        <span className="px-2 py-1 bg-emerald-500/20 text-emerald-400 text-[10px] font-bold rounded uppercase tracking-widest border border-emerald-500/30">
          Live Auction
        </span>
      </div>
      
      <div className="agent-list space-y-4">
        {sortedAgents.map(([agent, score], index) => (
          <div 
            key={agent} 
            className="agent-card group relative flex justify-between items-center p-3 bg-gradient-to-br from-white/5 to-transparent rounded-xl border border-white/5 hover:border-emerald-500/40 transition-all duration-300"
          >
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-full bg-slate-800 flex items-center justify-center border border-white/10 group-hover:scale-110 transition-transform">
                <span className="text-xs font-bold text-white/40">#{index + 1}</span>
              </div>
              <div>
                <div className="font-mono text-sm font-bold text-white/90 group-hover:text-emerald-400 transition-colors">{agent}</div>
                <div className="text-[10px] text-white/30 uppercase tracking-tighter">Senior Specialist</div>
              </div>
            </div>

            <div className="flex flex-col items-end gap-1">
              <div className="reputation-bar w-24 h-1.5 bg-white/5 rounded-full overflow-hidden">
                <div 
                  className="h-full bg-gradient-to-r from-emerald-600 to-emerald-400 transition-all duration-1000 ease-out" 
                  style={{ width: `${score * 100}%` }}
                />
              </div>
              <div className="flex items-center gap-1.5">
                <span className="text-lg font-black text-white leading-none">{(score * 10).toFixed(1)}</span>
                <span className="text-[10px] font-bold text-emerald-500/60 uppercase">CR</span>
              </div>
            </div>
            
            {/* Subtle glow effect on hover */}
            <div className="absolute inset-0 bg-emerald-500/5 opacity-0 group-hover:opacity-100 rounded-xl transition-opacity pointer-events-none" />
          </div>
        ))}
        {sortedAgents.length === 0 && (
          <div className="py-10 text-center text-white/20 italic text-sm">Waiting for agents to join the agora...</div>
        )}
      </div>
    </div>
  );
};

const ShieldView = ({ blockedActions }) => {
  return (
    <div className="shield-container glass-morphism p-6 rounded-2xl shadow-[0_20px_50px_rgba(239,68,68,0.1)] border border-red-500/10 bg-slate-900/40 backdrop-blur-xl">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-2xl font-black text-red-500 tracking-tight flex items-center gap-2">
          <span>🛡️ Immune Shield</span>
        </h3>
        <div className="flex gap-1">
          <div className="w-2 h-2 rounded-full bg-red-500 animate-pulse" />
          <div className="w-2 h-2 rounded-full bg-red-500/40" />
        </div>
      </div>

      <div className="action-logs space-y-3">
        {blockedActions.map((action, i) => (
          <div 
            key={i} 
            className="action-blocked group p-3 bg-red-950/20 border-l-2 border-red-500/60 rounded-lg font-mono text-xs text-red-200/80 hover:bg-red-950/40 transition-colors"
          >
            <div className="flex justify-between mb-1">
              <span className="text-red-500 font-bold uppercase text-[9px]">Critical Policy Violation</span>
              <span className="text-white/20">0.003ms</span>
            </div>
            <div className="break-all">{action}</div>
            <div className="mt-2 text-[9px] text-white/40 italic">Triggered by: internal/agent/laws.l</div>
          </div>
        ))}
        {blockedActions.length === 0 && (
          <div className="flex flex-col items-center py-8 gap-4">
            <div className="w-12 h-12 rounded-full bg-emerald-500/10 flex items-center justify-center border border-emerald-500/20">
              <div className="w-6 h-6 rounded-full border-2 border-emerald-500/40 border-t-emerald-500 animate-spin" />
            </div>
            <p className="text-white/30 text-[10px] uppercase tracking-widest font-bold">Scanning for anomalies...</p>
          </div>
        )}
      </div>
    </div>
  );
};

export { MarketplaceView, ShieldView };
