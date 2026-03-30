package telegram

import (
	"context"
	"time"

	"github.com/kocar/aurelia/internal/agent"
)

// ChatStreamToken represents a single streaming token or completion event
// sent to the dashboard chat UI.
type ChatStreamToken struct {
	Content      string `json:"content,omitempty"`
	Done         bool   `json:"done,omitempty"`
	Error        string `json:"error,omitempty"`
	InputTokens  int    `json:"input_tokens,omitempty"`
	OutputTokens int    `json:"output_tokens,omitempty"`
}

// StreamChat processes a user message through the full pipeline and streams
// tokens back via the returned channel. This is the web-dashboard equivalent
// of ProcessExternalInput but with token-level streaming.
func (bc *BotController) StreamChat(ctx context.Context, userID int64, text string) <-chan ChatStreamToken {
	out := make(chan ChatStreamToken, 100)

	go func() {
		defer close(out)

		ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
		defer cancel()

		session := bc.newInputSessionWithContext(ctx, userID, text)
		session.text = text

		// Porteiro security check
		if bc.porteiro != nil {
			safe, err := bc.porteiro.IsSafe(ctx, text)
			if err == nil && !safe {
				out <- ChatStreamToken{Error: "Mensagem bloqueada pelo sistema de segurança.", Done: true}
				return
			}
		}

		// Memory compression
		if bc.memory != nil {
			compressCtx, compressCancel := context.WithTimeout(ctx, 5*time.Second)
			_ = bc.memory.Compress(compressCtx, session.convID)
			compressCancel()
		}

		// Persist incoming message
		_ = bc.persistIncomingContext(session, userID)

		// Prepare execution context
		activeSkill, history, systemPrompt, allowedTools, err := bc.prepareExecution(session)
		if err != nil {
			out <- ChatStreamToken{Error: err.Error(), Done: true}
			return
		}

		// Stream from executor
		stream, err := bc.executor.ExecuteStream(ctx, systemPrompt, activeSkill, history, allowedTools)
		if err != nil {
			out <- ChatStreamToken{Error: err.Error(), Done: true}
			return
		}

		var fullContent string
		for resp := range stream {
			if resp.Err != nil {
				if fullContent != "" {
					// Partial content available — send it before the error
					out <- ChatStreamToken{Done: true, Content: fullContent}
				} else {
					out <- ChatStreamToken{Error: resp.Err.Error(), Done: true}
				}
				return
			}
			if resp.Content != "" {
				fullContent += resp.Content
				out <- ChatStreamToken{Content: resp.Content}
			}
			if resp.Done {
				out <- ChatStreamToken{
					Done:         true,
					InputTokens:  resp.InputTokens,
					OutputTokens: resp.OutputTokens,
				}
			}
		}

		// Persist assistant response
		if fullContent != "" {
			finalAnswer := sanitizeAssistantOutputForUser(fullContent)
			if bc.porteiro != nil {
				finalAnswer = bc.porteiro.SecureOutput(finalAnswer)
			}
			bc.persistAssistantAnswer(session, finalAnswer)
		}
	}()

	return out
}

// ChatHistory returns recent messages for the given user, suitable for
// rendering the chat history on the dashboard.
func (bc *BotController) ChatHistory(ctx context.Context, userID int64) ([]agent.Message, error) {
	session := bc.newInputSessionWithContext(ctx, userID, "")
	return bc.buildAgentHistory(session)
}
