//go:build integration

package dashboard

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freePort: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	return port
}

func TestDashboard_SSE_PublishAndReceive(t *testing.T) {
	port := freePort(t)
	logger := slog.Default()

	if err := StartServer(logger, port); err != nil {
		t.Fatalf("StartServer() error = %v", err)
	}
	// Aguarda o servidor subir
	time.Sleep(50 * time.Millisecond)

	received := make(chan string, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Conectar SSE client
	go func() {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
			fmt.Sprintf("http://127.0.0.1:%d/api/events", port), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data:") {
				select {
				case received <- strings.TrimPrefix(line, "data: "):
				default:
				}
				return
			}
		}
	}()

	// Dar tempo pro SSE client conectar
	time.Sleep(100 * time.Millisecond)

	// Publicar evento
	Publish(Event{
		Type:      "integration_test",
		Agent:     "smoke",
		Action:    "verify",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})

	select {
	case data := <-received:
		if !strings.Contains(data, "integration_test") {
			t.Fatalf("evento recebido não contém 'integration_test': %q", data)
		}
		t.Logf("SSE recebido OK: %q", data)
	case <-ctx.Done():
		t.Fatal("timeout: nenhum evento SSE recebido em 5s")
	}
}

func TestDashboard_SubscriberCleanup(t *testing.T) {
	before := subscriberCount()

	port := freePort(t)
	if err := StartServer(slog.Default(), port); err != nil {
		t.Fatalf("StartServer() error = %v", err)
	}
	time.Sleep(50 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		defer close(done)
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
			fmt.Sprintf("http://127.0.0.1:%d/api/events", port), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		bufio.NewScanner(resp.Body).Scan()
	}()

	time.Sleep(100 * time.Millisecond)
	if subscriberCount() <= before {
		t.Log("aviso: subscriber pode não ter registrado a tempo (race window)")
	}

	cancel()
	<-done
	time.Sleep(100 * time.Millisecond)

	if subscriberCount() != before {
		t.Fatalf("subscriber não foi removido após disconnect: got %d esperava %d",
			subscriberCount(), before)
	}
	t.Logf("subscriber cleanup OK: count=%d", subscriberCount())
}

func subscriberCount() int {
	subMu.Lock()
	defer subMu.Unlock()
	return len(subscribers)
}
