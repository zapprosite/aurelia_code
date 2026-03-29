// Package voice provides conversation context for the Jarvis voice loop
// ADR: 20260328-e2e-jarvis-loop-wake-tts

package voice

import (
	"sync"
	"time"
)

// Message represents a single message in the conversation
type Message struct {
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ConversationContext maintains conversation history for voice interactions
// ADR: 20260328-e2e-jarvis-loop-wake-tts
type ConversationContext struct {
	UserID       int64
	ChatID       int64
	History      []Message
	LastIntent   string
	LastEntities map[string]string
	TTL          time.Duration
	LastActivity time.Time
	mu           sync.RWMutex
}

const defaultContextTTL = 5 * time.Minute

// NewConversationContext creates a new conversation context
func NewConversationContext(userID, chatID int64) *ConversationContext {
	return &ConversationContext{
		UserID:       userID,
		ChatID:       chatID,
		History:      make([]Message, 0, 100),
		LastEntities: make(map[string]string),
		TTL:          defaultContextTTL,
		LastActivity: time.Now(),
	}
}

// Add adds a user and assistant message pair to the history
func (c *ConversationContext) Add(user, assistant string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.History = append(c.History, Message{
		Role:      "user",
		Content:   user,
		Timestamp: time.Now(),
	})
	c.History = append(c.History, Message{
		Role:      "assistant",
		Content:   assistant,
		Timestamp: time.Now(),
	})
	c.LastActivity = time.Now()

	// Prune if too long (keep last 50 exchanges)
	if len(c.History) > 100 {
		c.History = c.History[len(c.History)-100:]
	}
}

// IsExpired returns true if the context has expired
func (c *ConversationContext) IsExpired() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Since(c.LastActivity) > c.TTL
}

// GetHistory returns the conversation history as messages
func (c *ConversationContext) GetHistory() []Message {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]Message, len(c.History))
	copy(result, c.History)
	return result
}

// Clear resets the conversation history
func (c *ConversationContext) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.History = make([]Message, 0, 100)
	c.LastIntent = ""
	c.LastEntities = make(map[string]string)
}

// ContextStore manages conversation contexts
// ADR: 20260328-e2e-jarvis-loop-wake-tts
type ContextStore struct {
	contexts map[int64]*ConversationContext // key: userID
	mu       sync.RWMutex
}

// NewContextStore creates a new context store
func NewContextStore() *ContextStore {
	return &ContextStore{
		contexts: make(map[int64]*ConversationContext),
	}
}

// Get retrieves or creates a conversation context
func (s *ContextStore) Get(userID, chatID int64) *ConversationContext {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := userID // simplified - use chatID if needed
	ctx, exists := s.contexts[key]
	if !exists {
		ctx = NewConversationContext(userID, chatID)
		s.contexts[key] = ctx
	}

	// Reset TTL on access
	ctx.LastActivity = time.Now()
	return ctx
}

// Set updates an existing conversation context
func (s *ContextStore) Set(userID int64, ctx *ConversationContext) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.contexts[userID] = ctx
}

// Cleanup removes expired contexts
func (s *ContextStore) Cleanup() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	removed := 0
	for key, ctx := range s.contexts {
		if ctx.IsExpired() {
			delete(s.contexts, key)
			removed++
		}
	}
	return removed
}
