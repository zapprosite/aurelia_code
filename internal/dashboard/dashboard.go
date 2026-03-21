package dashboard

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
)

//go:embed dist
var content embed.FS

// StartServer inicia o servidor Web do Dashboard ULTRATRINK na porta 3333.
// Executa de forma assíncrona recebendo um logger.
func StartServer(logger *slog.Logger) error {
	subFS, err := fs.Sub(content, "dist")
	if err != nil {
		logger.Error("erro ao carregar arquivos estáticos do dashboard", slog.Any("err", err))
		return err
	}

	mux := http.NewServeMux()
	// Middleware simples para logs e servir o SPA
	mux.Handle("/", http.FileServer(http.FS(subFS)))

	go func() {
		logger.Info("ULTRATRINK Dashboard Online", slog.String("url", "http://localhost:3333"))
		if err := http.ListenAndServe(":3333", mux); err != nil {
			logger.Error("servidor do dashboard parou", slog.Any("err", err))
		}
	}()
	
	return nil
}
