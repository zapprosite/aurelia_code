import os
import json
import sqlite3
import requests
import sys

# Paths
PROJECT_ROOT = "/home/will/aurelia"
SOURCE_DIR = f"{PROJECT_ROOT}/.agent/skills"
DB_PATH = "/home/will/.aurelia/data/aurelia.db"
QDRANT_URL = "http://localhost:6333/collections/aurelia_skills"
OLLAMA_URL = "http://localhost:11434/api/embeddings"
EMBED_MODEL = "nomic-embed-text"

def get_embedding(text):
    try:
        response = requests.post(OLLAMA_URL, json={
            "model": EMBED_MODEL,
            "prompt": text
        }, timeout=30)
        return response.json().get("embedding")
    except Exception as e:
        print(f"Error getting embedding: {e}")
        return None

def main():
    print("🚀 Iniciando Indexador Semântico de Skills (SOTA 2026.1)...")
    
    # 1. SQLite Setup
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()
    cursor.execute("DROP TABLE IF EXISTS skills_meta")
    cursor.execute("""
        CREATE TABLE IF NOT EXISTS skills_meta (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT UNIQUE,
            path TEXT,
            description TEXT,
            last_sync TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            source TEXT DEFAULT '.agent'
        )
    """)
    conn.commit()

    # 2. Scan Skills
    skills = []
    for skill_dir in os.listdir(SOURCE_DIR):
        full_path = os.path.join(SOURCE_DIR, skill_dir)
        if os.path.isdir(full_path):
            skill_md = os.path.join(full_path, "SKILL.md")
            if os.path.exists(skill_md):
                with open(skill_md, "r") as f:
                    content = f.read()
                    # Simple extraction of first 200 chars for metadata
                    desc = content.split("\n")[2:10] if "---" in content else content[:200]
                    desc = " ".join(desc).strip()[:200]
                    skills.append({
                        "name": skill_dir,
                        "content": content,
                        "desc": desc,
                        "path": full_path
                    })

    # 3. Indexing Loop
    for i, skill in enumerate(skills):
        print(f"[{i+1}/{len(skills)}] Indexando {skill['name']}...")
        
        # A. Embedding
        vector = get_embedding(skill["content"])
        if not vector: continue

        # B. Qdrant Upsert
        point = {
            "points": [
                {
                    "id": i + 1000, # Offset to avoid collisions
                    "vector": vector,
                    "payload": {
                        "name": skill["name"],
                        "description": skill["desc"],
                        "path": skill["path"],
                        "source": ".agent"
                    }
                }
            ]
        }
        res = requests.put(f"{QDRANT_URL}/points?wait=true", json=point)
        
        # C. SQLite Update
        cursor.execute("""
            INSERT OR REPLACE INTO skills_meta (name, path, description, source)
            VALUES (?, ?, ?, ?)
        """, (skill["name"], skill["path"], skill["desc"], ".agent"))
        conn.commit()

    print(f"✅ Sincronização e Indexação concluídas. {len(skills)} habilidades processadas.")
    conn.close()

if __name__ == "__main__":
    main()
