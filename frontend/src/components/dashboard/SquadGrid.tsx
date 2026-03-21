import { useEffect, useState } from "react";
import * as Icons from "lucide-react";
import { Card } from "../ui/Card";
import { Badge } from "../ui/Badge";
import { motion } from "framer-motion";

export type SquadAgent = {
  id: string;
  name: string;
  role: string;
  status: string;
  load: number;
  color: string;
  icon: string;
};

const getIcon = (iconName: string) => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const Icon = (Icons as Record<string, any>)[iconName] || Icons.CircleDashed;
  return Icon;
};

export function SquadGrid() {
  const [agents, setAgents] = useState<SquadAgent[]>([]);

  const fetchSquad = () => {
    fetch('/api/squad')
      .then(res => res.json())
      .then(data => {
        if (data && Array.isArray(data)) {
          setAgents(data);
        }
      })
      .catch(err => console.error("Failed to fetch squad:", err));
  };

  useEffect(() => {
    fetchSquad();
    
    // Listen for SSE updates if they are broadcasted
    const eventSource = new EventSource("/api/events");
    eventSource.onmessage = (event) => {
       try {
         const data = JSON.parse(event.data);
         if (data.type === 'agent_status_update' || data.type === 'agent_tool' || data.type === 'agent_handoff') {
            // Se um agente executou algo ou foi atualizado, vamos dar fetch no squad inteiro 
            // (ou aplicar o diff se houver payload especifico).
            // Fetch é seguro e rápido localmente.
            fetchSquad();
         }
       } catch {
         // ignore invalid json from sse
       }
    };

    return () => eventSource.close();
  }, []);

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      {agents.map((agent, i) => {
        const Icon = getIcon(agent.icon);
        return (
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
                    <Icon className="w-5 h-5" />
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
        );
      })}
    </div>
  );
}
