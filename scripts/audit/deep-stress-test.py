import requests
import time
import json

def test_tier(label, prompt, timeout=None):
    url = "http://localhost:4000/chat/completions"
    headers = {
        "Content-Type": "application/json",
        "Authorization": "Bearer sk-1234"
    }
    data = {
        "model": "aurelia-smart",
        "messages": [{"role": "user", "content": prompt}],
        "max_tokens": 50
    }
    
    # Se quisermos testar o timeout do T0 (10s), podemos enviar uma carga pesada
    # Ou simplesmente ver o que o roteador escolhe
    
    print(f"\n🔍 [AUDITORIA] {label}")
    start = time.time()
    try:
        response = requests.post(url, headers=headers, json=data)
        duration = time.time() - start
        
        if response.status_code == 200:
            res_json = response.json()
            # O LiteLLM às vezes retorna o modelo real no campo 'model' da resposta
            # dependendo da configuração. Se não, o log do container mostra.
            actual_model = res_json.get('model', 'unknown')
            content = res_json['choices'][0]['message']['content'].strip()[:100]
            
            # Tentar pegar o provedor nos headers customizados (se ativado no litellm)
            provider = response.headers.get('x-litellm-model-id', 'aurelia-smart-group')
            
            print(f"✅ Sucesso em {duration:.2f}s")
            print(f"📦 Modelo Real: {actual_model}")
            print(f"📡 Provedor Headers: {provider}")
            print(f"💬 Resposta: {content}...")
        else:
            print(f"❌ Erro {response.status_code}: {response.text}")
    except Exception as e:
        print(f"💥 Falha Crítica: {e}")

print("=== 🧪 AURELIA SMART DEEP AUDIT (SOTA 2026.1) ===")

# Teste 1: Rapidez Local (T0)
test_tier("NÍVEL 0: Gemma 3 Local", "Diga 'Local'.")

# Teste 2: Groq Turbo (T1) -> Vamos pedir algo que o T0 pode falhar ou demorar
# (Simulação: Vamos rodar 5 em paralelo para estressar o T0 se for o caso)
test_tier("NÍVEL 1: Cascata de Nuvem", "Escreva um código em Go para um servidor gRPC de alto desempenho.")

print("\n=== ✅ FIM DA AUDITORIA DE ESTRESSE ===")
