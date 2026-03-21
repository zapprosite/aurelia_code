package dashboard

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestUltratrinkE2E(t *testing.T) {
	// Preparação do handler do logger
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	// Dado que o servidor do ULTRATRINK está integrado e serve os arquivos do Vite
	err := StartServer(logger)
	if err != nil {
		t.Fatalf("falha severa ao instanciar servidor em memória: %v", err)
	}

	// Aguarda a rotina HTTP do ListenAndServe acoplar a porta de teste
	time.Sleep(1 * time.Second)

	// Quando os Agentes em enxame emitirem status ou a visão do humano tentar acessar 
	// o Painel na porta :3333...
	resp, err := http.Get("http://localhost:3333/")
	if err != nil {
		t.Fatalf("erro ao simular acesso do cliente E2E: %v", err)
	}
	defer resp.Body.Close()

	// Então o servidor deve retornar HTTP 200 VÁLIDO com o index.html gerado pelo React
	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperado HTTP 200 OK, recebeu %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil || len(bodyBytes) == 0 {
		t.Errorf("expected HTML payload from single page app, got empty or error: %v", err)
	}
	
	t.Logf("E2E OK. Dashboard do enxame subiu servindo %d bytes de estáticos e pronto pro 24/7.", len(bodyBytes))
}
