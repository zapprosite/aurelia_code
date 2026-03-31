#!/usr/bin/env python3
import socket
import os
import sys
import time
import numpy as np

# Configurações Industriais SOTA 2026
CHUNK = 1024
CHANNELS = 1
RATE = 16000
SOCKET_PATH = "/tmp/aurelia-voice.sock"

def main():
    print("🛰️ Voice Gateway online (MOCK MODE - No PyAudio required)")

    # Limpa socket antigo
    if os.path.exists(SOCKET_PATH):
        os.remove(SOCKET_PATH)

    server = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    server.bind(SOCKET_PATH)
    server.listen(1)
    
    print(f"🛰️ Voice Gateway online em {SOCKET_PATH}")

    while True:
        conn, addr = server.accept()
        print("🔗 Go Actor conectado ao stream de áudio.")
        try:
            # Injeta 1s de ruído para teste SOTA 2026.1
            noise = np.random.randint(-20000, 20000, RATE, dtype=np.int16)
            conn.sendall(noise.tobytes())
            
            while True:
                # Simula leitura de áudio (silêncio ou ruído baixo)
                data = np.random.randint(-100, 100, CHUNK, dtype=np.int16).tobytes()
                conn.sendall(data)
                time.sleep(CHUNK / RATE)
        except (ConnectionResetError, BrokenPipeError):
            print("❌ Go Actor desconectado.")
        finally:
            conn.close()

if __name__ == "__main__":
    main()
