package a2a

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/observability"
)

//go:generate stringer -type=MessageRole -output=messages_string.go

// MessageRole is the role of the sender in an A2A message.
type MessageRole string

const (
	RoleUser    MessageRole = "user"
	RoleAgent   MessageRole = "agent"
	RoleSystem  MessageRole = "system"
)

// Message is the core unit of A2A communication.
type Message struct {
	// Unique message ID.
	ID string `json:"id"`
	// Role of the sender.
	Role MessageRole `json:"role"`
	// Message content parts.
	Parts []MessagePart `json:"parts"`
	// Timestamp of the message.
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// MessagePart is the content of a message. Exactly one field is set.
type MessagePart interface {
	partMarker()
}

// TextPart contains plain text.
type TextPart struct {
	Text string `json:"text"`
}

func (TextPart) partMarker() {}

// DataPart contains structured JSON data.
type DataPart struct {
	Data any `json:"data"`
}

func (DataPart) partMarker() {}

// SendMessageRequest is the A2A SendMessage request body.
type SendMessageRequest struct {
	Message   Message  `json:"message"`
	Stream   bool      `json:"stream,omitempty"`
	TaskID   string   `json:"taskId,omitempty"`
	SessionID string   `json:"sessionId,omitempty"`
}

// SendMessageResponse is the A2A SendMessage response body.
type SendMessageResponse struct {
	Message Message `json:"message"`
	TaskID string  `json:"taskId,omitempty"`
}

// MessageProcessor processes incoming A2A messages and produces responses.
type MessageProcessor interface {
	ProcessMessage(ctx context.Context, msg Message, opts ProcessOptions) (*MessageProcessingResult, error)
}

// ProcessOptions provides context for message processing.
type ProcessOptions struct {
	SessionID string
	TaskID    string
}

// MessageProcessingResult is the result of processing a message.
type MessageProcessingResult struct {
	// Content of the response.
	Content []MessagePart
	// Whether this is a streaming response.
	Streaming bool
	// Optional task status.
	TaskStatus *TaskStatus
}

// TaskStatus describes the state of a long-running task.
type TaskStatus struct {
	State   TaskState `json:"state"`
	Message string    `json:"message,omitempty"`
}

// TaskState is the state of a task.
type TaskState string

const (
	TaskStatePending   TaskState = "pending"
	TaskStateWorking  TaskState = "working"
	TaskStateComplete TaskState = "completed"
	TaskStateFailed   TaskState = "failed"
	TaskStateCancel   TaskState = "canceled"
)

// extractText returns the first text part from a message, or empty string.
func extractText(msg Message) string {
	for _, part := range msg.Parts {
		if tp, ok := part.(TextPart); ok {
			return tp.Text
		}
	}
	return ""
}

// UnmarshalMessagePart deserializes a MessagePart from JSON.
func UnmarshalMessagePart(raw json.RawMessage) (MessagePart, error) {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, err
	}
	if _, ok := obj["text"]; ok {
		var tp TextPart
		if err := json.Unmarshal(raw, &tp); err != nil {
			return nil, err
		}
		return tp, nil
	}
	if _, ok := obj["data"]; ok {
		var dp DataPart
		if err := json.Unmarshal(raw, &dp); err != nil {
			return nil, err
		}
		return dp, nil
	}
	return nil, fmt.Errorf("unknown message part type")
}

// UnmarshalMessage deserializes a Message from JSON.
func UnmarshalMessage(raw json.RawMessage) (*Message, error) {
	var m struct {
		ID        string          `json:"id"`
		Role      MessageRole      `json:"role"`
		Parts     json.RawMessage `json:"parts"`
		Timestamp *time.Time      `json:"timestamp,omitempty"`
	}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	parts := make([]MessagePart, 0)
	for _, rawPart := range parseArray(raw, "parts") {
		part, err := UnmarshalMessagePart(rawPart)
		if err != nil {
			continue
		}
		parts = append(parts, part)
	}
	msg := &Message{ID: m.ID, Role: m.Role, Parts: parts}
	if m.Timestamp != nil {
		msg.Timestamp = *m.Timestamp
	}
	return msg, nil
}

func parseArray(raw json.RawMessage, key string) []json.RawMessage {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil
	}
	arr, ok := obj[key]
	if !ok {
		return nil
	}
	var result []json.RawMessage
	json.Unmarshal(arr, &result)
	return result
}

// A2AServer handles A2A protocol HTTP requests.
type A2AServer struct {
	card       *AgentCard
	processor MessageProcessor
	log       *slog.Logger
	authToken string
}

// NewA2AServer creates a new A2A server.
func NewA2AServer(card *AgentCard, processor MessageProcessor, authToken string) *A2AServer {
	return &A2AServer{
		card:       card,
		processor: processor,
		log:       observability.Logger("a2a.server"),
		authToken: authToken,
	}
}

// ServeHTTP implements http.Handler.
func (s *A2AServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	switch path {
	case ".well-known/agent.json":
		s.serveAgentCard(w, r)
		return
	case "a2a":
		s.serveA2A(w, r)
		return
	default:
		http.NotFound(w, r)
	}
}

// serveAgentCard returns the AgentCard at /.well-known/agent.json.
func (s *A2AServer) serveAgentCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(s.card)
}

// serveA2A handles the A2A SendMessage endpoint.
func (s *A2AServer) serveA2A(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Bearer token auth
	if s.authToken != "" {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") || strings.TrimPrefix(auth, "Bearer ") != s.authToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1MB max
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var req SendMessageRequest
	if err := json.Unmarshal(body, &req); err != nil {
		s.log.Warn("invalid A2A request body", slog.Any("err", err))
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	opts := ProcessOptions{
		SessionID: req.SessionID,
		TaskID:    req.TaskID,
	}

	result, err := s.processor.ProcessMessage(ctx, req.Message, opts)
	if err != nil {
		s.log.Error("A2A message processing failed", slog.Any("err", err))
		// Return error as text part
		result = &MessageProcessingResult{
			Content: []MessagePart{TextPart{Text: fmt.Sprintf("erro: %v", err)}},
		}
	}

	// Build response
	resp := SendMessageResponse{
		Message: Message{
			ID:    fmt.Sprintf("resp-%d", time.Now().UnixNano()),
			Role:  RoleAgent,
			Parts: result.Content,
		},
	}

	if result.TaskStatus != nil {
		resp.Message.ID = req.Message.ID // same task
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
