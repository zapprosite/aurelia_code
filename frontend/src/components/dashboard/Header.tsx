import { Cpu, Search, Bell } from "lucide-react";
import { Badge } from "../ui/Badge";
import { Button } from "../ui/Button";

interface HeaderProps {
  title: string;
}

export function Header({ title }: HeaderProps) {
  return (
    <header className="h-16 flex items-center justify-between px-8 border-b border-white/10 bg-white/5 backdrop-blur-sm z-40 sticky top-0">
      <div className="flex items-center gap-4">
        <h2 className="text-xl font-medium tracking-tight text-white/90">{title}</h2>
        <Badge variant="outline" className="opacity-50 text-[10px] py-0">Main Floor</Badge>
      </div>

      <div className="flex items-center gap-6">
        {/* System Stats Component (Quick View) */}
        <div className="hidden md:flex items-center gap-4 px-4 py-1.5 rounded-full bg-white/5 border border-white/5 text-[11px] font-mono text-white/40">
          <div className="flex items-center gap-2">
            <Cpu className="w-3 h-3 text-primary/60" />
            <span>VRAM: <span className="text-white/70">8.2 GiB</span></span>
          </div>
          <div className="h-3 w-px bg-white/10" />
          <div className="flex items-center gap-2">
             <div className="w-1.5 h-1.5 rounded-full bg-blue-400 shadow-[0_0_5px_rgba(96,165,250,0.5)]" />
             <span>GPU: <span className="text-white/70">32%</span></span>
          </div>
        </div>

        <div className="flex items-center gap-2">
           <Button variant="ghost" size="icon" className="text-white/40 hover:text-white">
              <Search className="w-4 h-4" />
           </Button>
           <Button variant="ghost" size="icon" className="text-white/40 hover:text-white relative">
              <Bell className="w-4 h-4" />
              <span className="absolute top-2 right-2 w-1.5 h-1.5 bg-red-500 rounded-full border border-[#161616]" />
           </Button>
        </div>
      </div>
    </header>
  );
}
