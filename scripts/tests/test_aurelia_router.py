import litellm
import time
import os

# Configuração do endpoint local do seu roteador LiteLLM
LITELLM_PROXY_URL = os.environ.get("LITELLM_LOCAL_URL", "http://localhost:4000")
MASTER_KEY = os.environ.get("LITELLM_MASTER_KEY", "{{LITELLM_MASTER_KEY_HIDDEN}}b83cfa00b8b8c90c966367cf1bee8efdc5b504c3b05872011d3085ebc737a3cb")

def test_tier_routing(prompt, description):
    print(f"\n--- Teste: {description} ---")
    start_time = time.time()
    
    try:
        response = litellm.completion(
            model="openai/aurelia-smart", # O nome que definimos no config.yaml
            messages=[{"role": "user", "content": prompt}],
            api_base=LITELLM_PROXY_URL,
            api_key=MASTER_KEY,
            timeout=35 # Tempo para permitir a cascata completa
        )
        
        duration = time.time() - start_time
        model_used = response.get("_source_model", "Desconhecido") # LiteLLM injeta isso no log
        content = response.choices[0].message.content[:100]
        
        print(f"✅ Sucesso em {duration:.2f}s")
        print(f"🧠 Modelo que respondeu: {model_used}")
        print(f"📝 Resposta (resumo): {content}...")
        
    except Exception as e:
        print(f"❌ Falha no teste: {str(e)}")

# --- CENÁRIO 1: Soberania Local (Tier 0) ---
# Esperado: Resposta rápida do Gemma 3 na sua 4090
test_tier_routing(
    "Responda apenas com 'LOCAL OK' se você for o Gemma 3 rodando na RTX 4090.",
    "Soberania Local (Tier 0)"
)

# --- CENÁRIO 2: Simulação de Carga/Complexidade (Tier 1) ---
# Esperado: Se o prompt for longo ou o local falhar, o Minimax ou Gemini assume.
test_tier_routing(
    "Explique detalhadamente a teoria das cordas e como ela se integra à gravidade quântica em 3 parágrafos.",
    "Salto para Nuvem Gratuita (Tier 1)"
)

# --- CENÁRIO 3: Teste de Embeddings (Nomic Local) ---
print("\n--- Teste: Embedding Local (Nomic) ---")
try:
    embed_resp = litellm.embedding(
        model="openai/nomic-embed-text",
        input=["Aurelia Smart é o futuro da IA soberana."],
        api_base=LITELLM_PROXY_URL,
        api_key=MASTER_KEY
    )
    print(f"✅ Embedding gerado com sucesso! Dimensões: {len(embed_resp['data'][0]['embedding'])}")
except Exception as e:
    print(f"❌ Falha no Embedding: {str(e)}")
