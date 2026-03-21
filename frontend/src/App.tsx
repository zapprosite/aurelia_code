import { Terminal, Users, Cpu, FileBox, LayoutList, Settings, Sparkles } from 'lucide-react';
import { motion } from 'framer-motion';

function App() {
  return (
    <div className="flex h-screen bg-[#0f0f0f] text-gray-100 font-sans overflow-hidden">
      {/* Sidebar - Squad Navigation */}
      <nav className="w-64 border-r border-white/10 bg-[#161616]/80 backdrop-blur-xl flex flex-col">
        <div className="h-16 flex items-center px-6 border-b border-white/10">
          <Sparkles className="w-5 h-5 text-blue-500 mr-2" />
          <h1 className="text-sm font-bold tracking-widest text-white/90">ULTRATRINK</h1>
        </div>
        
        <div className="flex-1 py-6 flex flex-col gap-2">
          <NavItem icon={<ActivityIcon />} label="Timeline & Feed" active />
          <NavItem icon={<Users />} label="Squad (Agentes)" />
          <NavItem icon={<FileBox />} label="The Brain (.context)" />
          <NavItem icon={<LayoutList />} label="Roadmap Slices" />
        </div>

        <div className="p-4 border-t border-white/10 text-xs text-white/40 flex items-center justify-between">
          <span>Aurelia Daemon</span>
          <span className="flex items-center gap-1">
            <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
            Stable v5.0
          </span>
        </div>
      </nav>

      {/* Main Content Area */}
      <main className="flex-1 flex flex-col">
        {/* Top Header */}
        <header className="h-16 flex items-center justify-between px-8 border-b border-white/10 bg-white/5 backdrop-blur-sm">
          <h2 className="text-xl font-medium tracking-tight text-white/90">Main Floor — Activity Stream</h2>
          <div className="flex items-center gap-4">
            <Cpu className="w-4 h-4 text-white/40" />
            <span className="text-xs font-mono text-white/50">VRAM: 8.2 GiB / 24 GiB (RX 7900 XTX)</span>
            <Settings className="w-4 h-4 text-white/40 cursor-pointer hover:text-white transition-colors" />
          </div>
        </header>

        {/* Dashboard Feed Grid */}
        <div className="flex-1 p-8 overflow-y-auto">
          <div className="max-w-5xl mx-auto space-y-6">
            
            {/* Feed Component Mock */}
            <motion.div 
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              className="p-6 rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-2xl"
            >
              <div className="flex items-center gap-3 mb-4">
                <div className="w-8 h-8 rounded-full bg-blue-500/20 flex items-center justify-center border border-blue-500/30">
                  <Terminal className="w-4 h-4 text-blue-400" />
                </div>
                <div>
                  <h3 className="text-sm font-medium text-white/90">Antigravity (Orquestrador)</h3>
                  <p className="text-xs text-white/40">commit docs(adr): incluir roadmap do Dashboard ULTRATRINK</p>
                </div>
                <span className="ml-auto text-xs text-white/30 font-mono">Just Now</span>
              </div>
              
              <div className="bg-[#0a0a0a] rounded-xl p-4 border border-white/5 font-mono text-xs text-emerald-400/80">
                <p>{'>'} git push --no-verify origin HEAD</p>
                <p className="text-white/40 mt-1">To https://github.com/zapprosite/aurelia_code.git</p>
                <p className="text-white/40">   b95ce85..5c580c7  HEAD {'->'} feat/gpu-cpu-rate-limit</p>
              </div>
            </motion.div>

            {/* Second Mock */}
            <motion.div 
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 }}
              className="p-6 rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-2xl opacity-60"
            >
              <div className="flex items-center gap-3 mb-4">
                <div className="w-8 h-8 rounded-full bg-purple-500/20 flex items-center justify-center border border-purple-500/30">
                  <span className="text-xs font-bold text-purple-400">A</span>
                </div>
                <div>
                  <h3 className="text-sm font-medium text-white/90">Aurelia (Base)</h3>
                  <p className="text-xs text-white/40">Gerenciando estado da memória sync</p>
                </div>
                <span className="ml-auto text-xs text-white/30 font-mono">2 mins ago</span>
              </div>
              <p className="text-sm text-white/70">Qdrant re-index pipeline completed successfully for vector db #memory-sync.</p>
            </motion.div>

          </div>
        </div>
      </main>
    </div>
  );
}

// Subcomponents
function ActivityIcon() {
  return (
    <svg className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M22 12h-4l-3 9L9 3l-3 9H2" />
    </svg>
  );
}

function NavItem({ icon, label, active = false }: { icon: React.ReactNode, label: string, active?: boolean }) {
  return (
    <button 
      className={`
        w-full flex items-center gap-3 px-6 py-2.5 text-sm transition-all duration-200
        ${active 
          ? 'text-white bg-white/10 border-r-2 border-blue-500' 
          : 'text-white/50 hover:text-white/90 hover:bg-white/5'
        }
      `}
    >
      <span className={active ? 'text-blue-400' : 'text-white/40'}>{icon}</span>
      <span className="font-medium tracking-wide">{label}</span>
    </button>
  );
}

export default App;
