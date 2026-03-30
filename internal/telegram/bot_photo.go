package telegram

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"sort"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/observability"
	"gopkg.in/telebot.v3"
)

func (bc *BotController) handlePhoto(c telebot.Context) error {
	m := c.Message()
	if m.AlbumID != "" {
		if bc.storeAlbumPhoto(m.AlbumID, m.ID, m.Caption, *m.Photo) {
			go func(albumID string) {
				time.Sleep(2 * time.Second)
				bc.processAlbum(c, albumID)
			}(m.AlbumID)
		}
		return nil
	}

	return bc.processPhoto(c, m.Caption, *m.Photo)
}

func (bc *BotController) processPhoto(c telebot.Context, caption string, photo telebot.Photo) error {
	stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.Typing, 4*time.Second)
	defer stopTyping()

	filePath, err := bc.downloadTelegramFile(&photo.File, photo.FileID+".jpg")
	if err != nil {
		observability.Logger("telegram.vision").Warn("failed to download photo", slog.Any("err", err))
		return SendContextText(c, downloadFailureMessage)
	}

	data, err := os.ReadFile(filePath)

	if err != nil {
		return SendContextText(c, "Erro ao ler imagem baixada.")
	}

	mimeType := detectImageMIMEType(filePath, "")
	parts := []agent.ContentPart{
		{Type: agent.ContentPartText, Text: caption},
		{Type: agent.ContentPartImage, MIMEType: mimeType, Data: data},
	}

	return bc.processInputWithParts(c, parts)
}

func (bc *BotController) processAlbum(c telebot.Context, albumID string) {
	caption, photos, ok := bc.flushAlbumPhotos(albumID)
	if !ok {
		return
	}

	parts := []agent.ContentPart{{Type: agent.ContentPartText, Text: caption}}
	for _, p := range photos {
		filePath, err := bc.downloadTelegramFile(&p.photo.File, p.photo.FileID+".jpg")
		if err != nil {
			continue
		}
		data, err := os.ReadFile(filePath)

		if err != nil {
			continue
		}
		parts = append(parts, agent.ContentPart{
			Type:     agent.ContentPartImage,
			MIMEType: detectImageMIMEType(filePath, ""),
			Data:     data,
		})
	}

	if err := bc.processInputWithParts(c, parts); err != nil {
		observability.Logger("telegram.vision").Warn("album processing failed", slog.Any("err", err))
		_ = SendError(bc.bot, c.Chat(), "Falha ao processar o álbum de fotos.")
	}
}

func (bc *BotController) storeAlbumPhoto(albumID string, messageID int, caption string, photo telebot.Photo) bool {
	bc.albumMu.Lock()
	defer bc.albumMu.Unlock()

	album, exists := bc.pendingAlbums[albumID]
	if !exists {
		bc.pendingAlbums[albumID] = &pendingAlbum{
			ownerMessageID: messageID,
			caption:        caption,
			photos:         []albumPhoto{{messageID: messageID, photo: photo}},
		}
		return true
	}

	album.photos = append(album.photos, albumPhoto{messageID: messageID, photo: photo})
	if caption != "" {
		album.caption = caption
	}
	return false
}

func (bc *BotController) flushAlbumPhotos(albumID string) (string, []albumPhoto, bool) {
	bc.albumMu.Lock()
	defer bc.albumMu.Unlock()

	album, exists := bc.pendingAlbums[albumID]
	if !exists {
		return "", nil, false
	}

	delete(bc.pendingAlbums, albumID)
	sort.Slice(album.photos, func(i, j int) bool {
		return album.photos[i].messageID < album.photos[j].messageID
	})

	return album.caption, album.photos, true
}

func isSupportedImageDocument(filename, mimeType string) bool {
	mimeType = strings.ToLower(mimeType)
	if strings.HasPrefix(mimeType, "image/") {
		return true
	}
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return true
	}
	return false
}

func detectImageMIMEType(filename, mimeType string) string {
	if strings.HasPrefix(strings.ToLower(mimeType), "image/") {
		return mimeType
	}
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return "image/jpeg"
	}
}

func (bc *BotController) storeRecentMedia(session inputSession) {
	if len(session.message.Parts) <= 1 {
		return
	}

	bc.mediaMu.Lock()
	defer bc.mediaMu.Unlock()

	var mediaParts []agent.ContentPart
	for _, p := range session.message.Parts {
		if p.Type == agent.ContentPartImage {
			mediaParts = append(mediaParts, p)
		}
	}

	if len(mediaParts) > 0 {
		bc.recentMedia[session.convID] = recentMedia{
			parts:     mediaParts,
			updatedAt: time.Now(),
		}
	}
}

func (bc *BotController) loadRecentMedia(convID string) (recentMedia, bool) {
	bc.mediaMu.Lock()
	defer bc.mediaMu.Unlock()

	media, exists := bc.recentMedia[convID]
	if !exists {
		return recentMedia{}, false
	}

	if time.Since(media.updatedAt) > 30*time.Minute {
		delete(bc.recentMedia, convID)
		return recentMedia{}, false
	}

	return media, true
}

func shouldReuseRecentMedia(text string) bool {
	low := strings.ToLower(text)
	triggers := []string{"imagem", "foto", "anexo", "pdf", "mandado", "referencia"}
	for _, t := range triggers {
		if strings.Contains(low, t) {
			return true
		}
	}
	return false
}

func (bc *BotController) attachRecentMediaIfRelevant(c telebot.Context, text string, parts []agent.ContentPart) []agent.ContentPart {
	if !shouldReuseRecentMedia(text) {
		return parts
	}

	media, ok := bc.loadRecentMedia(fmt.Sprintf("%d", c.Sender().ID))
	if !ok {
		return parts
	}

	if parts == nil {
		parts = []agent.ContentPart{{Type: agent.ContentPartText, Text: text}}
	}

	parts = append(parts, media.parts...)
	return parts
}

// Helpers missing from previous listing but required
func (bc *BotController) processInputWithParts(c telebot.Context, parts []agent.ContentPart) error {
	// Re-wrap input session for multimodal
	session := bc.newInputSession(c, "")
	session.message.Parts = parts
	session.text = session.persistedContent()

	bc.storeRecentMedia(session)
	// processInput as usual
	return bc.processInputSession(c, session, false)
}

func (s *inputSession) persistedContent() string {
	images := 0
	caption := ""
	for _, p := range s.message.Parts {
		if p.Type == agent.ContentPartImage {
			images++
		}
		if p.Type == agent.ContentPartText && p.Text != "" {
			caption = p.Text
		}
	}
	if images > 0 {
		desc := fmt.Sprintf("[%d imagem(ns) enviada(s)]", images)
		if caption != "" {
			return fmt.Sprintf("%s\n%s", caption, desc)
		}
		return desc
	}
	return s.text
}
