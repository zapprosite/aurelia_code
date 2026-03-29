import * as React from "react";
import { cn } from "../../lib/utils";
import { Sparkles, Activity, Users, Brain, Map, Settings, ChevronLeft, ChevronRight, Server, Bot, Monitor } from "lucide-react";
import { Button } from "../ui/Button";
import { motion } from "framer-motion";

export type TabId = "timeline" | "bots" | "squad" | "brain" | "roadmap" | "homelab" | "computer";

interface SidebarProps {
  activeTab: TabId;
  setActiveTab: (tab: TabId) => void;
}

export function Sidebar({ activeTab, setActiveTab }: SidebarProps) {
  const [collapsed, setCollapsed] = React.useState(false);

  const navItems = [
    { id: "timeline", label: "Timeline & Feed", icon: Activity },
    { id: "bots",     label: "Team Bots",        icon: Bot },
    { id: "computer", label: "Computer Use",    icon: Monitor },
    { id: "squad",    label: "Squad (Agentes)",   icon: Users },
    { id: "brain",    label: "The Brain (.context)", icon: Brain },
    { id: "roadmap",  label: "Roadmap Slices",    icon: Map },
    { id: "homelab",  label: "Homelab",           icon: Server },
  ];

  return (
    <motion.nav 
      initial={false}
      animate={{ width: collapsed ? 80 : 256 }}
      className={cn(
        "relative flex flex-col h-screen border-r border-white/10 bg-[#161616]/80 backdrop-blur-xl transition-all duration-300 ease-in-out",
        collapsed ? "items-center" : ""
      )}
    >
      {/* Brand Header */}
      <div className={cn(
        "h-16 flex items-center px-6 border-b border-white/10",
        collapsed ? "justify-center px-0" : "justify-between"
      )}>
        {!collapsed && (
          <div className="flex items-center gap-2 overflow-hidden truncate">
            <Sparkles className="w-5 h-5 text-primary shrink-0" />
            <span className="text-sm font-bold tracking-widest text-white/90">ULTRATRINK</span>
          </div>
        )}
        {collapsed && <Sparkles className="w-6 h-6 text-primary" />}
        
        <Button 
          variant="ghost" 
          size="icon" 
          className="absolute -right-3 top-20 z-50 h-6 w-6 rounded-full border border-white/10 bg-[#161616] hover:bg-white/5 shadow-xl"
          onClick={() => setCollapsed(!collapsed)}
        >
          {collapsed ? <ChevronRight className="h-3 w-3" /> : <ChevronLeft className="h-3 w-3" />}
        </Button>
      </div>

      {/* Navigation Main */}
      <div className="flex-1 py-8 flex flex-col gap-1 px-3">
        {navItems.map((item) => (
          <NavItem 
            key={item.id}
            icon={<item.icon />}
            label={item.label}
            active={activeTab === item.id}
            collapsed={collapsed}
            onClick={() => setActiveTab(item.id as TabId)}
          />
        ))}
      </div>

      {/* Footer Status */}
      <div className={cn(
        "p-4 border-t border-white/10 bg-white/5 flex flex-col gap-2",
        collapsed ? "items-center" : ""
      )}>
        <div className={cn("flex items-center gap-3", collapsed ? "justify-center" : "justify-between")}>
          {!collapsed && <span className="text-[10px] font-medium text-white/40 uppercase tracking-tighter">Aurelia v5.0</span>}
          <div className="flex items-center gap-1.5 uppercase">
             <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </span>
            {!collapsed && <span className="text-[10px] font-bold text-emerald-400/90">Stable</span>}
          </div>
        </div>
        {!collapsed && <Settings className="w-4 h-4 text-white/20 mt-2 cursor-pointer hover:text-white/60 transition-colors" />}
      </div>
    </motion.nav>
  );
}

function NavItem({ icon, label, active, collapsed, onClick }: { 
  icon: React.ReactNode, 
  label: string, 
  active: boolean, 
  collapsed: boolean,
  onClick: () => void 
}) {
  return (
    <button 
      onClick={onClick}
      className={cn(
        "group relative flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all duration-200 overflow-hidden",
        active 
          ? "bg-primary/10 text-primary shadow-[inset_0_0_1px_1px_rgba(59,130,246,0.1)]" 
          : "text-white/40 hover:text-white/90 hover:bg-white/5 active:scale-95"
      )}
    >
      <div className={cn(
        "transition-colors duration-200",
        active ? "text-primary shadow-[0_0_10px_rgba(59,130,246,0.5)]" : "text-white/40 group-hover:text-white/60"
      )}>
        {icon}
      </div>
      
      {!collapsed && (
        <span className={cn(
          "font-medium tracking-wide text-sm whitespace-nowrap transition-all",
          active ? "text-white/90" : ""
        )}>
          {label}
        </span>
      )}

      {active && (
        <motion.div 
          layoutId="sidebar-active-indicator"
          className="absolute left-0 w-1 h-2/3 bg-primary rounded-r-full"
        />
      )}
    </button>
  );
}
