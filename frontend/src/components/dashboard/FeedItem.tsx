import { Terminal, Github, Container, Sparkles, Repeat, Volume2 } from "lucide-react";
import { Card } from "../ui/Card";
import { Badge } from "../ui/Badge";
import { motion } from "framer-motion";

export type FeedType = "git" | "docker" | "ai" | "system" | "handoff";

export interface FeedItemProps {
  id: string;
  type: FeedType;
  agent: string;
  action: string;
  timestamp: string;
  content?: string;
  status?: "pending" | "success" | "warning" | "error";
}

export function FeedItem({ type, agent, action, timestamp, content, status = "success" }: FeedItemProps) {
  const getIcon = () => {
    switch (type) {
      case "git": return <Github className="w-4 h-4 text-blue-400" />;
      case "docker": return <Container className="w-4 h-4 text-blue-300" />;
      case "ai": return <Sparkles className="w-4 h-4 text-purple-400" />;
      case "handoff": return <Repeat className="w-4 h-4 text-orange-400" />;
      default: return <Terminal className="w-4 h-4 text-emerald-400" />;
    }
  };

  const getAgentColor = () => {
    switch (agent) {
      case "Antigravity": return "bg-blue-500/20 border-blue-500/30 text-blue-400";
      case "Aurelia": return "bg-purple-500/20 border-purple-500/30 text-purple-400";
      case "Claude": return "bg-orange-500/20 border-orange-500/30 text-orange-400";
      default: return "bg-white/10 border-white/20 text-white/60";
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 15 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, scale: 0.95 }}
      layout
    >
      <Card className="p-5 border-white/5 hover:border-white/10 transition-colors group">
        <div className="flex items-center gap-4 mb-4">
          <div className={`w-10 h-10 rounded-xl flex items-center justify-center border ${getAgentColor()}`}>
            {getIcon()}
          </div>
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2">
               <h3 className="text-sm font-bold text-white/90 uppercase tracking-wide">{agent}</h3>
               <Badge variant="outline" className="text-[9px] uppercase tracking-tighter opacity-40">{type}</Badge>
            </div>
            <p className="text-xs text-white/50 truncate font-medium mt-0.5">{action}</p>
          </div>
          <div className="text-right shrink-0">
             <span className="text-[10px] text-white/20 font-mono block">{timestamp}</span>
             <div className="flex justify-end items-center gap-2 mt-1">
                {type === "ai" && (
                  <button className="p-1 hover:bg-white/10 rounded-md transition-colors group/audio">
                    <Volume2 className="w-3 h-3 text-white/20 group-hover/audio:text-purple-400 transition-colors" />
                  </button>
                )}
                <span className={`w-1.5 h-1.5 rounded-full ${status === "success" ? "bg-emerald-500 shadow-[0_0_5px_rgba(16,185,129,0.5)]" : "bg-orange-500"}`} />
             </div>
          </div>
        </div>
        
        {content && (
          <div className="bg-[#0a0a0a]/80 rounded-xl p-4 border border-white/5 font-mono text-xs leading-relaxed text-blue-300/80 shadow-inner group-hover:border-white/10 transition-colors">
            {content.split("\n").map((line, i) => (
              <p key={i} className={line.startsWith(">") ? "text-emerald-400/90" : ""}>{line}</p>
            ))}
          </div>
        )}
      </Card>
    </motion.div>
  );
}
