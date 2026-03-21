import { Shield, Brain, Zap, Terminal } from "lucide-react";
import { Card } from "../ui/Card";
import { Badge } from "../ui/Badge";
import { motion } from "framer-motion";

const agents = [
  { id: "aurelia", name: "Aurélia", role: "Arquiteta / Governança", status: "online", load: 12, icon: Shield, color: "text-purple-400" },
  { id: "antigravity", name: "Antigravity", role: "Orquestrador / Cockpit", status: "busy", load: 68, icon: Zap, color: "text-blue-400" },
  { id: "claude", name: "Claude 5 Omni", role: "Implementador Principal", status: "online", load: 5, icon: Brain, color: "text-orange-400" },
  { id: "codex", name: "Codex", role: "Executor Rápido", status: "offline", load: 0, icon: Terminal, color: "text-emerald-400" },
];

export function SquadGrid() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      {agents.map((agent, i) => (
        <motion.div
           key={agent.id}
           initial={{ opacity: 0, scale: 0.9 }}
           animate={{ opacity: 1, scale: 1 }}
           transition={{ delay: i * 0.1 }}
        >
          <Card className="p-4 border-white/5 hover:border-white/20 transition-all cursor-default overflow-hidden relative group">
            <div className={`absolute top-0 right-0 w-24 h-24 blur-3xl opacity-10 transition-opacity group-hover:opacity-20 ${agent.color.replace("text-", "bg-")}`} />
            
            <div className="flex items-center justify-between mb-4 relative z-10">
               <div className={`p-2 rounded-lg bg-white/5 border border-white/5 ${agent.color}`}>
                  <agent.icon className="w-5 h-5" />
               </div>
               <Badge variant={agent.status === 'online' ? 'success' : agent.status === 'busy' ? 'default' : 'outline'} className="uppercase text-[9px]">
                 {agent.status}
               </Badge>
            </div>

            <div className="relative z-10">
              <h4 className="font-bold text-white/90 tracking-tight">{agent.name}</h4>
              <p className="text-[11px] text-white/40 font-medium mb-3">{agent.role}</p>

              <div className="space-y-1.5">
                <div className="flex justify-between text-[10px] text-white/30 font-mono">
                  <span>Workloadpool</span>
                  <span>{agent.load}%</span>
                </div>
                <div className="h-1 w-full bg-white/5 rounded-full overflow-hidden">
                   <motion.div 
                     initial={{ width: 0 }}
                     animate={{ width: `${agent.load}%` }}
                     className={`h-full ${agent.color.replace("text-", "bg-")}`} 
                   />
                </div>
              </div>
            </div>
          </Card>
        </motion.div>
      ))}
    </div>
  );
}
