// Package vision provides screen capturing and state tracking for the agent.
package vision

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

// CacheEntry represents a cached vision result or state.
type CacheEntry struct {
	ImageHash string
	State     *ScreenState
	Timestamp time.Time
}

// VisionCache implements a semantic cache for visual states to avoid redundant LLM calls.
type VisionCache struct {
	entries map[string]CacheEntry
	mu      sync.RWMutex
	maxSize int
}

// NewVisionCache creates a new vision cache with a specified maximum size.
func NewVisionCache(maxSize int) *VisionCache {
	return &VisionCache{
		entries: make(map[string]CacheEntry),
		maxSize: maxSize,
	}
}

// GetHash computes a SHA256 hash of the image data.
func (c *VisionCache) GetHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// Store adds a new entry to the cache.
func (c *VisionCache) Store(imageHash string, state *ScreenState) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Simple eviction: clear if too big (industrial 2026 approach)
	if len(c.entries) >= c.maxSize {
		c.entries = make(map[string]CacheEntry)
	}

	c.entries[imageHash] = CacheEntry{
		ImageHash: imageHash,
		State:     state,
		Timestamp: time.Now(),
	}
}

// Get retrieves a cached state by its image hash.
func (c *VisionCache) Get(imageHash string) (*ScreenState, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[imageHash]
	if !ok {
		return nil, false
	}

	// Invalidate after 5 minutes (standard for dynamic web content)
	if time.Since(entry.Timestamp) > 5*time.Minute {
		return nil, false
	}

	return entry.State, true
}
