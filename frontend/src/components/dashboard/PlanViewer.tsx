import { CheckCircle2, Circle, Clock, AlertTriangle, ShieldCheck, XCircle } from "lucide-react";
import { motion } from "framer-motion";
import { Card } from "../ui/Card";
import { Badge } from "../ui/Badge";

export interface PlanStep {
  order: number;
  description: string;
  tool?: string;
  status: "pending" | "done";
}

export interface ActionPlan {
  id: string;
  title: string;
  description: string;
  risk_level: "low" | "medium" | "high" | "critical";
  steps: PlanStep[];
  estimated_time: string;
  backout_plan: string;
  created_at: string;
  status: "proposed" | "approved" | "rejected" | "executing" | "completed";
}

interface PlanViewerProps {
  plans: ActionPlan[];
  onAction: (planId: string, action: "approve" | "reject") => void;
}

export function PlanViewer({ plans, onAction }: PlanViewerProps) {
  if (plans.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-32 text-center opacity-40">
        <Clock className="w-12 h-12 mb-4 text-white/20" />
        <p className="text-sm font-mono uppercase tracking-widest text-white/40">No active plans in queue</p>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {plans.map((plan) => (
        <motion.div
          key={plan.id}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="relative"
        >
          <Card className={`overflow-hidden border-l-4 ${
            plan.risk_level === 'critical' ? 'border-l-red-500' :
            plan.risk_level === 'high' ? 'border-l-orange-500' :
            plan.risk_level === 'medium' ? 'border-l-yellow-500' : 'border-l-blue-500'
          }`}>
            <div className="p-6">
              <div className="flex items-start justify-between mb-6">
                <div>
                  <div className="flex items-center gap-3 mb-2">
                    <h3 className="text-xl font-bold text-white/90">{plan.title}</h3>
                    <Badge variant={
                      plan.status === 'approved' ? 'success' : 
                      plan.status === 'rejected' ? 'destructive' : 'default'
                    } className="uppercase text-[10px]">
                      {plan.status}
                    </Badge>
                  </div>
                  <p className="text-sm text-white/40 max-w-2xl">{plan.description}</p>
                </div>
                <div className="text-right">
                  <div className="text-[10px] text-white/20 font-mono uppercase tracking-widest mb-1">Risk Assessment</div>
                  <div className={`flex items-center gap-1.5 justify-end font-bold uppercase text-xs ${
                    plan.risk_level === 'critical' ? 'text-red-400' :
                    plan.risk_level === 'high' ? 'text-orange-400' :
                    plan.risk_level === 'medium' ? 'text-yellow-400' : 'text-blue-400'
                  }`}>
                    {plan.risk_level === 'critical' && <AlertTriangle className="w-3.5 h-3.5" />}
                    {plan.risk_level}
                  </div>
                </div>
              </div>

              <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div className="lg:col-span-2 space-y-3">
                  <div className="text-[10px] text-white/20 font-mono uppercase tracking-widest">Execution Steps</div>
                  <div className="space-y-2">
                    {plan.steps.map((step) => (
                      <div key={step.order} className="flex items-center gap-4 p-3 rounded-lg bg-white/5 border border-white/5 group hover:border-white/10 transition-all">
                        <div className="flex-shrink-0">
                          {step.status === 'done' ? (
                            <CheckCircle2 className="w-4 h-4 text-green-500" />
                          ) : (
                            <Circle className="w-4 h-4 text-white/20" />
                          )}
                        </div>
                        <div className="flex-1">
                          <p className="text-sm text-white/70">{step.description}</p>
                          {step.tool && <code className="text-[10px] text-primary/60 font-mono">{step.tool}</code>}
                        </div>
                        <div className="text-[10px] text-white/10 font-mono italic">#{step.order}</div>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="space-y-6">
                  <div>
                    <div className="text-[10px] text-white/20 font-mono uppercase tracking-widest mb-3">Resources & Safety</div>
                    <div className="space-y-3">
                      <div className="flex items-center gap-3 text-xs text-white/60">
                        <Clock className="w-4 h-4 text-white/20" />
                        <span>EST: <b className="text-white/80">{plan.estimated_time}</b></span>
                      </div>
                      <div className="flex items-start gap-3 text-xs text-white/60">
                        <ShieldCheck className="w-4 h-4 text-white/20 flex-shrink-0" />
                        <div>
                          <div className="text-white/30 text-[9px] uppercase tracking-tighter mb-0.5">Backout Strategy</div>
                          <span className="leading-relaxed italic">{plan.backout_plan}</span>
                        </div>
                      </div>
                    </div>
                  </div>

                  {plan.status === 'proposed' && (
                    <div className="pt-4 flex gap-3">
                      <button 
                        onClick={() => onAction(plan.id, "approve")}
                        className="flex-1 flex items-center justify-center gap-2 py-2.5 rounded-xl bg-green-500/10 hover:bg-green-500/20 text-green-500 border border-green-500/20 text-xs font-bold transition-all"
                      >
                        <ShieldCheck className="w-4 h-4" />
                        APPROVE
                      </button>
                      <button 
                        onClick={() => onAction(plan.id, "reject")}
                        className="flex-1 flex items-center justify-center gap-2 py-2.5 rounded-xl bg-red-500/10 hover:bg-red-500/20 text-red-500 border border-red-500/20 text-xs font-bold transition-all"
                      >
                        <XCircle className="w-4 h-4" />
                        REJECT
                      </button>
                    </div>
                  )}
                </div>
              </div>
            </div>
            
            {plan.status === 'executing' && (
              <div className="h-1 w-full bg-white/5 relative overflow-hidden">
                <motion.div 
                  initial={{ x: "-100%" }}
                  animate={{ x: "100%" }}
                  transition={{ repeat: Infinity, duration: 1.5, ease: "linear" }}
                  className="absolute inset-0 w-1/3 bg-gradient-to-r from-transparent via-primary/40 to-transparent shadow-[0_0_10px_rgba(var(--primary),0.5)]" 
                />
              </div>
            )}
          </Card>
        </motion.div>
      ))}
    </div>
  );
}
