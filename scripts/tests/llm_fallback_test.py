import litellm
import time
import os
import sys

# Configuração do endpoint local do seu roteador LiteLLM
LITELLM_PROXY_URL = os.environ.get("LITELLM_LOCAL_URL", "http://localhost:4000")
MASTER_KEY = os.environ.get("LITELLM_MASTER_KEY", "sk-1234")

def test_tier_routing(prompt, description):
    print(f"\n--- Teste: {description} ---")
    start_time = time.time()
    
    try:
        # Nota: Usamos 'openai/aurelia-smart' para garantir compatibilidade via LiteLLM Proxy
        response = litellm.completion(
            model="openai/aurelia-smart",
            messages=[{"role": "user", "content": prompt}],
            api_base=LITELLM_PROXY_URL,
            api_key=MASTER_KEY,
            timeout=35
        )
        
        duration = time.time() - start_time
        # LiteLLM injeta o modelo real no campo '_source_model' ou 'model' da resposta
        model_used = response.get("model", "Desconhecido")
        content = response.choices[0].message.content[:150]
        
        print(f"✅ Sucesso em {duration:.2f}s")
        print(f"🧠 Modelo Ativo: {model_used}")
        print(f"📝 Resposta: {content}...")
        
    except Exception as e:
        print(f"❌ Falha no teste: {str(e)}")

def run_suite():
    print("🚀 Iniciando Bateria de Testes: Soberania Híbrida (v2026.1)")

    # --- CENÁRIO 1: Soberania Local (Tier 0) ---
    test_tier_routing(
        "Responda apenas com 'LOCAL OK' se você for um modelo local operando via Ollama.",
        "Soberania Local (Tier 0)"
    )

    # --- CENÁRIO 2: Simulação de Complexidade (Tier 1/2) ---
    test_tier_routing(
        "Analise a arquitetura de microserviços em Rust vs Go sob a ótica de latência zero e segurança de memória. 3 parágrafos técnicos.",
        "Cascata de Inteligência (T1/T2)"
    )

    # --- CENÁRIO 3: Teste de Embeddings ---
    print("\n--- Teste: Embedding Local (Nomic) ---")
    try:
        embed_resp = litellm.embedding(
            model="openai/nomic-embed-text",
            input=["Aurelia Smart: Soberania local, escalabilidade global."],
            api_base=LITELLM_PROXY_URL,
            api_key=MASTER_KEY
        )
        print(f"✅ Embedding gerado! Dimensões: {len(embed_resp['data'][0]['embedding'])}")
    except Exception as e:
        print(f"❌ Falha no Embedding: {str(e)}")

if __name__ == "__main__":
    run_suite()
