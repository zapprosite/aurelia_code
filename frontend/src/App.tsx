import * as React from "react";
import { AnimatePresence, motion } from "framer-motion";
import { Sidebar, type TabId } from "./components/sidebar/Sidebar";
import { Header } from "./components/dashboard/Header";
import { FeedItem, type FeedItemProps } from "./components/dashboard/FeedItem";
import { AgentDesktop } from "./components/dashboard/AgentDesktop";
import { AgentComms } from "./components/dashboard/AgentComms";
import { CommandMenu } from "./components/dashboard/CommandMenu";
import { PlanViewer, type ActionPlan } from "./components/dashboard/PlanViewer";
import { ScrollArea } from "./components/ui/ScrollArea";
import { Card, CardHeader, CardTitle } from "./components/ui/Card";
import { Badge } from "./components/ui/Badge";
import { Layout, Server } from "lucide-react";
import { BrainSearch } from "./components/dashboard/BrainSearch";
import { BotsTab } from "./components/dashboard/BotsTab";
import { useSystemMetrics } from "./hooks/useSystemMetrics";
import { type SquadAgent } from "./components/dashboard/SquadGrid";

// Initial Feed Placeholder
const INITIAL_FEED: FeedItemProps[] = [
  {
    id: "welcome",
    type: "system",
    agent: "System",
    action: "Awaiting Live Events...",
    timestamp: "Now",
    content: "Connected to ULTRATRINK Real-time Engine. Monitoring Aurelia's neural activity.",
    status: "success"
  }
];

// S-28: VRV Homelab tab with manual refresh button
function HomelabTab() {
  const iframeRef = React.useRef<HTMLIFrameElement>(null);
  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2 text-white/40">
        <Server className="w-4 h-4" />
        <span className="text-xs font-mono uppercase tracking-widest">VRV — Homelab Monitor</span>
        <button
          onClick={() => {
            if (iframeRef.current) iframeRef.current.src = iframeRef.current.src;
          }}
          className="ml-auto font-mono text-[11px] px-2 py-0.5 rounded border border-white/10 text-white/30 hover:text-white/60 hover:border-white/30 transition-colors"
          title="Atualizar painel VRV"
        >
          [↺ atualizar]
        </button>
      </div>
      <div className="rounded-xl overflow-hidden border border-white/10" style={{ height: "75vh" }}>
        <iframe
          ref={iframeRef}
          src="/api/vrv/"
          className="w-full h-full"
          style={{ border: "none", background: "#0a0a0a" }}
          title="VRV Homelab Dashboard"
        />
      </div>
    </div>
  );
}

function App() {
  const [activeTab, setActiveTab] = React.useState<TabId>("timeline");
  const [feed, setFeed] = React.useState<FeedItemProps[]>(INITIAL_FEED);
  const [plans, setPlans] = React.useState<ActionPlan[]>([]);
  const [squadAgents, setSquadAgents] = React.useState<SquadAgent[]>([]);
  const metrics = useSystemMetrics();

  React.useEffect(() => {
    // SSE Connection to Go Backend
    const eventSource = new EventSource("/api/events");

    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        const newItem: FeedItemProps = {
          id: Math.random().toString(36).substr(2, 9),
          type: data.type === "agent_thought" ? "ai" :
                (data.type === "agent_tool" ? "system" :
                (data.type === "agent_handoff" ? "handoff" : "git")),
          agent: data.agent,
          action: data.action,
          timestamp: data.timestamp || "Just Now",
          content: typeof data.payload === "string" ? data.payload : JSON.stringify(data.payload, null, 2),
          status: "success"
        };

        if (data.type === "agent_plan") {
           setPlans(prev => {
              const existing = prev.find(p => p.id === data.payload.id);
              if (existing) {
                return prev.map(p => p.id === data.payload.id ? data.payload : p);
              }
              return [data.payload, ...prev];
           });
           if (data.action.includes("Proposto")) {
              setFeed(prev => [newItem, ...prev].slice(0, 50));
           }
        } else {
           setFeed(prev => [newItem, ...prev].slice(0, 50));
        }

        // Refresh squad on relevant events
        if (
          data.type === "agent_status_update" ||
          data.type === "agent_tool" ||
          data.type === "agent_handoff"
        ) {
          fetchSquad();
        }
      } catch (err) {
        console.error("Error parsing SSE event:", err);
      }
    };

    eventSource.onerror = (err) => {
      console.error("SSE connection error:", err);
    };

    return () => {
      eventSource.close();
    };
  }, []);

  const fetchSquad = () => {
    fetch("/api/squad")
      .then((res) => res.json())
      .then((data) => {
        if (Array.isArray(data)) setSquadAgents(data);
      })
      .catch((err) => console.error("Failed to fetch squad:", err));
  };

  React.useEffect(() => {
    fetchSquad();
  }, []);

  const getTabTitle = () => {
    switch (activeTab) {
      case "timeline": return "Main Floor — Activity Stream";
      case "bots":     return "Team Bots — Multi-Bot Manager";
      case "squad":    return "Squad View — Agent Status";
      case "brain":    return "The Brain — Semantic Context";
      case "roadmap":  return "Master Plan — Feature Roadmap";
      case "homelab":  return "Homelab — VRV Dashboard";
      default: return "Dashboard";
    }
  };

  // Derived squad stats
  const onlineCount = squadAgents.filter(a => a.status === "online" || a.status === "busy").length;
  const avgLoad = squadAgents.length > 0
    ? Math.round(squadAgents.reduce((s, a) => s + a.load, 0) / squadAgents.length)
    : 0;
  const gpuDisplay = metrics.gpuUtil !== null ? `${Math.round(metrics.gpuUtil)}%` : "N/A";

  return (
    <div className="flex h-screen bg-background text-foreground font-sans overflow-hidden selection:bg-primary/30">
      <CommandMenu />
      <Sidebar activeTab={activeTab} setActiveTab={setActiveTab} />

      <main className="flex-1 flex flex-col relative overflow-hidden">
        <div className="absolute top-[-10%] right-[-10%] w-[50%] h-[50%] bg-primary/5 blur-[120px] rounded-full pointer-events-none" />
        <div className="absolute bottom-[-10%] left-[-10%] w-[40%] h-[40%] bg-purple-500/5 blur-[120px] rounded-full pointer-events-none" />

        <Header title={getTabTitle()} />

        <ScrollArea className="flex-1">
          <div className="p-8 max-w-6xl mx-auto w-full">
            <AnimatePresence mode="wait">
              <motion.div
                key={activeTab}
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ duration: 0.2, ease: "easeOut" }}
                className="min-h-full"
              >
                {activeTab === "bots" && (
                  <BotsTab feed={feed} />
                )}

                {activeTab === "timeline" && (
                  <div className="space-y-6">
                     <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center gap-2 text-white/40">
                           <Layout className="w-4 h-4" />
                           <span className="text-xs font-mono uppercase tracking-widest">Real-time Operations Log</span>
                        </div>
                        <Badge variant="outline" className="text-[10px] opacity-30">AUTO-SYNC: ON</Badge>
                     </div>
                     <div className="space-y-4">
                        {feed.map((item) => (
                           <FeedItem key={item.id} {...item} />
                        ))}
                     </div>
                  </div>
                )}

                {activeTab === "squad" && (
                  <div className="space-y-8">
                     <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                        <Card className="bg-primary/5 border-primary/20">
                           <CardHeader className="pb-2">
                              <CardTitle className="text-xs uppercase text-primary tracking-widest">Agentes Online</CardTitle>
                              <div className="text-2xl font-bold text-white/90">{onlineCount}</div>
                           </CardHeader>
                        </Card>
                        <Card>
                           <CardHeader className="pb-2">
                              <CardTitle className="text-xs uppercase text-white/30 tracking-widest">Carga Média</CardTitle>
                              <div className="text-2xl font-bold text-white/90">{avgLoad}%</div>
                           </CardHeader>
                        </Card>
                        <Card>
                           <CardHeader className="pb-2">
                              <CardTitle className="text-xs uppercase text-white/30 tracking-widest">GPU</CardTitle>
                              <div className="text-2xl font-bold text-white/90">{gpuDisplay}</div>
                           </CardHeader>
                        </Card>
                     </div>
                     <AgentDesktop />
                     <AgentComms />
                  </div>
                )}

                {activeTab === "brain" && (
                  <BrainSearch />
                )}

                {activeTab === "homelab" && (
                  <HomelabTab />
                )}

                {activeTab === "roadmap" && (
                  <PlanViewer
                    plans={plans}
                    onAction={(planId, action) => {
                       fetch("/api/commands", {
                          method: "POST",
                          body: JSON.stringify({ action: `${action}_plan`, params: { plan_id: planId } })
                       });
                    }}
                  />
                )}
              </motion.div>
            </AnimatePresence>
          </div>
        </ScrollArea>
      </main>
    </div>
  );
}

export default App;
