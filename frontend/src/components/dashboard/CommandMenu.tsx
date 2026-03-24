import * as React from "react";
import { Search, Cpu, RefreshCw, Trash2 } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

interface CommandOption {
  id: string;
  label: string;
  description: string;
  icon: React.ElementType;
  action: string;
  params?: Record<string, string>;
}

const COMMANDS: CommandOption[] = [
  {
    id: "sync",
    label: "Sync Knowledge",
    description: "Re-index codebase documentation and memory",
    icon: RefreshCw,
    action: "sync_ai"
  },
  {
    id: "gemma",
    label: "Set Model: Gemma 3",
    description: "Switch primary LLM to Gemma 3 12B (Local)",
    icon: Cpu,
    action: "set_model",
    params: { model: "gemma3:12b" }
  },
  {
    id: "flush",
    label: "Flush Agent Memory",
    description: "Clear current conversation and task context",
    icon: Trash2,
    action: "flush_memory"
  }
];

export function CommandMenu() {
  const [isOpen, setIsOpen] = React.useState(false);
  const [search, setSearch] = React.useState("");

  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setIsOpen((open) => !open);
      }
    };

    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  const executeCommand = async (cmd: CommandOption) => {
    try {
      const resp = await fetch("/api/commands", {
        method: "POST",
        body: JSON.stringify({ action: cmd.action, params: cmd.params }),
      });
      if (resp.ok) {
        setIsOpen(false);
      }
    } catch (err) {
      console.error("Failed to execute command:", err);
    }
  };

  const filtered = COMMANDS.filter(c => 
    c.label.toLowerCase().includes(search.toLowerCase()) ||
    c.description.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={() => setIsOpen(false)}
            className="fixed inset-0 bg-black/60 backdrop-blur-sm z-[100]"
          />
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: -20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: -20 }}
            className="fixed top-[20%] left-1/2 -translate-x-1/2 w-full max-w-xl z-[101] px-4"
          >
            <div className="bg-[#121214]/80 backdrop-blur-2xl border border-white/10 rounded-2xl shadow-2xl overflow-hidden">
              <div className="flex items-center px-4 py-4 border-b border-white/5 gap-3">
                <Search className="w-5 h-5 text-white/20" />
                <input 
                  autoFocus
                  placeholder="Type a command or search..."
                  className="flex-1 bg-transparent border-none outline-none text-white/90 placeholder:text-white/20 text-lg"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                />
                <div className="flex items-center gap-1.5 px-2 py-1 bg-white/5 rounded border border-white/10 text-[10px] text-white/40 font-mono">
                  ESC
                </div>
              </div>

              <div className="max-h-[300px] overflow-y-auto p-2">
                {filtered.map((cmd) => (
                  <button
                    key={cmd.id}
                    onClick={() => executeCommand(cmd)}
                    className="w-full flex items-center gap-4 px-4 py-3 rounded-xl hover:bg-white/5 transition-colors text-left group"
                  >
                    <div className="w-10 h-10 rounded-lg bg-white/5 border border-white/5 flex items-center justify-center group-hover:border-primary/40 group-hover:text-primary transition-all">
                      <cmd.icon className="w-5 h-5 opacity-40 group-hover:opacity-100" />
                    </div>
                    <div>
                      <div className="text-sm font-bold text-white/80">{cmd.label}</div>
                      <div className="text-xs text-white/30">{cmd.description}</div>
                    </div>
                  </button>
                ))}
                {filtered.length === 0 && (
                  <div className="py-12 text-center text-white/20 text-sm">
                    No commands found matching "{search}"
                  </div>
                )}
              </div>

              <div className="px-4 py-3 bg-white/5 border-t border-white/5 flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <div className="flex items-center gap-1.5 text-[10px] text-white/30 uppercase tracking-widest">
                    <span>Select</span>
                    <div className="px-1 py-0.5 bg-white/10 rounded border border-white/10">↵</div>
                  </div>
                </div>
                <div className="text-[10px] text-white/20 font-mono uppercase tracking-widest">
                  Aurelia Cockpit v1.0
                </div>
              </div>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
}
