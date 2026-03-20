import React from 'react';
import DeckGL from '@deck.gl/react';
import { HeatmapLayer } from '@deck.gl/aggregation-layers';

// Mock de dados de "Atenção" dos agentes
const attentionData = [
  { position: [-46.6333, -23.5505], focus: 100, agent: "Regente" }, // São Paulo
  { position: [-46.6350, -23.5510], focus: 80, agent: "Dev-01" },
  { position: [-46.6370, -23.5520], focus: 60, agent: "Sec-01" },
];

export default function CognitiveAROverlay() {
  const layers = [
    new HeatmapLayer({
      id: 'attention-heatmap',
      data: attentionData,
      getPosition: d => d.position,
      getWeight: d => d.focus,
      radiusPixels: 60,
    })
  ];

  return (
    <div style={{ position: 'relative', width: '100%', height: '600px', backgroundColor: '#000' }}>
      <DeckGL
        initialViewState={{ longitude: -46.6333, latitude: -23.5505, zoom: 15 }}
        controller={true}
        layers={layers}
      />
      <div style={{ position: 'absolute', top: 20, left: 20, color: '#00ffcc', pointerEvents: 'none' }}>
        <h3>Cognitive AR Layer (2026)</h3>
        <p>Visualizando "balões de pensamento" e heatmaps de atividade mental sobre os dados.</p>
        <div style={{ backgroundColor: 'rgba(0,0,0,0.7)', padding: '10px', borderRadius: '5px' }}>
          <b>Atenção do Regente:</b> Analisando gargalo de latência na edge...
        </div>
      </div>
    </div>
  );
}
