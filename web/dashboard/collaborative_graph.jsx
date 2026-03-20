import React from 'react';
import ReactFlow, { Background, Controls } from 'reactflow';
import 'reactflow/dist/style.css';

const initialNodes = [
  { id: '1', data: { label: 'Conductor (Regente)' }, position: { x: 250, y: 0 }, style: { border: '2px solid #00ffcc', borderRadius: '8px', padding: '10px' } },
  { id: '2', data: { label: 'Specialist A (Dev)' }, position: { x: 100, y: 150 } },
  { id: '3', data: { label: 'Specialist B (Sec)' }, position: { x: 400, y: 150 } },
];

const initialEdges = [
  { id: 'e1-2', source: '1', target: '2', label: 'assign_task', animated: true },
  { id: 'e1-3', source: '1', target: '3', label: 'assign_task', animated: true },
  { id: 'e2-3', source: '2', target: '3', label: 'help_requested (SQL)', style: { stroke: '#ff0077' }, animated: true },
];

// Context Fusion Panel (Emblending 2026)
const ContextFusionPanel = () => (
  <div style={{ padding: '15px', backgroundColor: '#1a1a2e', border: '1px solid #00ffcc', borderRadius: '8px', color: '#fff', marginBottom: '10px' }}>
    <h4 style={{ margin: '0 0 10px 0' }}>🧠 Painel de Fusão Contextual (Emblending)</h4>
    <div style={{ display: 'flex', gap: '20px' }}>
      <div style={{ flex: 1 }}>
        <div style={{ fontSize: '12px' }}>Visão (VL) - 60%</div>
        <div style={{ height: '8px', width: '100%', backgroundColor: '#444', borderRadius: '4px' }}>
          <div style={{ height: '100%', width: '60%', backgroundColor: '#00ffcc', borderRadius: '4px' }}></div>
        </div>
      </div>
      <div style={{ flex: 1 }}>
        <div style={{ fontSize: '12px' }}>MCP (Legacy) - 30%</div>
        <div style={{ height: '8px', width: '100%', backgroundColor: '#444', borderRadius: '4px' }}>
          <div style={{ height: '100%', width: '30%', backgroundColor: '#ffcc00', borderRadius: '4px' }}></div>
        </div>
      </div>
      <div style={{ flex: 1 }}>
        <div style={{ fontSize: '12px' }}>Memória (RAG) - 10%</div>
        <div style={{ height: '8px', width: '100%', backgroundColor: '#444', borderRadius: '4px' }}>
          <div style={{ height: '100%', width: '10%', backgroundColor: '#ff0055', borderRadius: '4px' }}></div>
        </div>
      </div>
    </div>
    <p style={{ fontSize: '11px', marginTop: '8px', color: '#aaa' }}>Status: Grounding Multimodal Ativo - Zero Alucinações Detectadas</p>
  </div>
);

export default function SwarmOfficeGraph() {
  return (
    <div style={{ width: '100%', height: '100vh', display: 'flex', flexDirection: 'column', backgroundColor: '#0b0e14', padding: '20px' }}>
      <ContextFusionPanel />
      <div style={{ flex: 1, border: '1px solid #333', borderRadius: '8px' }}>
        <ReactFlow
          nodes={initialNodes}
          edges={initialEdges}
          fitView
        >
          <Background color="#1a1a1a" gap={16} />
          <Controls />
        </ReactFlow>
      </div>
    </div>
  );
}
