import requests
import time

url = "http://localhost:4000/chat/completions"
headers = {"Authorization": "Bearer sk-1234"}

models = [
    ("TIER 0 [Local]", "ollama/gemma3:12b"),
    ("TIER 1 [Minimax]", "openrouter/minimax/minimax-m2.5"),
    ("TIER 1 [Gemini]", "gemini/gemini-2.0-flash"),
    ("TIER 1 [Groq Turbo]", "groq/meta-llama/llama-4-scout-17b-16e-instruct"),
    ("TIER 2 [Expert]", "openrouter/moonshotai/kimi-k2.5")
]

print("=== 🚀 AURELIA SMART TIERS: VISUALIZAÇÃO SOTA 2026.1 ===")

for label, model_id in models:
    print(f"\n📡 {label}: {model_id}")
    data = {
        "model": model_id,
        "messages": [{"role": "user", "content": "Diga o seu nome e nível."}],
        "max_tokens": 50
    }
    
    start = time.time()
    try:
        response = requests.post(url, headers=headers, json=data)
        duration = time.time() - start
        
        if response.status_code == 200:
            content = response.json()['choices'][0]['message']['content'].strip()
            print(f"✅ OK em {duration:.2f}s | Resposta: {content}")
        else:
            print(f"❌ Erro {response.status_code}: {response.text}")
    except Exception as e:
        print(f"💥 Falha: {e}")

print("\n=== ✅ TESTE INDIVIDUAL CONCLUÍDO ===")
