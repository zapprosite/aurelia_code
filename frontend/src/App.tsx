import * as React from "react";
import { AnimatePresence, motion } from "framer-motion";
import { Sidebar, type TabId } from "./components/sidebar/Sidebar";
import { Header } from "./components/dashboard/Header";
import { FeedItem, type FeedItemProps } from "./components/dashboard/FeedItem";
import { SquadGrid } from "./components/dashboard/SquadGrid";
import { CommandMenu } from "./components/dashboard/CommandMenu";
import { PlanViewer, type ActionPlan } from "./components/dashboard/PlanViewer";
import { ScrollArea } from "./components/ui/ScrollArea";
import { Card, CardHeader, CardTitle } from "./components/ui/Card";
import { Badge } from "./components/ui/Badge";
import { Brain, Layout, type LucideIcon } from "lucide-react";

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

function App() {
  const [activeTab, setActiveTab] = React.useState<TabId>("timeline");
  const [feed, setFeed] = React.useState<FeedItemProps[]>(INITIAL_FEED);
  const [plans, setPlans] = React.useState<ActionPlan[]>([]);

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

  const getTabTitle = () => {
    switch (activeTab) {
      case "timeline": return "Main Floor — Activity Stream";
      case "squad": return "Squad View — Agent Status";
      case "brain": return "The Brain — Semantic Context";
      case "roadmap": return "Master Plan — Feature Roadmap";
      default: return "Dashboard";
    }
  };

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
                              <CardTitle className="text-xs uppercase text-primary tracking-widest">Avg Pulse Rate</CardTitle>
                              <div className="text-2xl font-bold text-white/90">98.2%</div>
                           </CardHeader>
                        </Card>
                        <Card>
                           <CardHeader className="pb-2">
                              <CardTitle className="text-xs uppercase text-white/30 tracking-widest">Active Threads</CardTitle>
                              <div className="text-2xl font-bold text-white/90">32</div>
                           </CardHeader>
                        </Card>
                        <Card>
                           <CardHeader className="pb-2">
                              <CardTitle className="text-xs uppercase text-white/30 tracking-widest">Memory Sync</CardTitle>
                              <div className="text-2xl font-bold text-white/90">OK</div>
                           </CardHeader>
                        </Card>
                     </div>
                     <SquadGrid />
                  </div>
                )}

                {activeTab === "brain" && <SectionPlaceholder icon={Brain} title="Semantic Cortex" description="Explorador de memória vetorial e documentos de contexto do projeto." />}
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

function SectionPlaceholder({ icon: Icon, title, description }: { icon: LucideIcon, title: string, description: string }) {
  return (
    <div className="flex flex-col items-center justify-center py-32 text-center">
       <div className="w-16 h-16 rounded-2xl bg-white/5 border border-white/5 flex items-center justify-center mb-6">
          <Icon className="w-8 h-8 text-white/20" />
       </div>
       <h3 className="text-xl font-bold text-white/80 mb-2 uppercase tracking-tight">{title}</h3>
       <p className="text-sm text-white/30 max-w-md">{description}</p>
       <Badge variant="outline" className="mt-6 uppercase text-[10px] tracking-widest">Integrating with Local Vector DB</Badge>
    </div>
  );
}

export default App;
