// Package vision provides screen state tracking for computer use
// ADR: 20260328-vision-pipeline-computer-use

package vision

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"
)

// ScreenState tracks the history of screen states during computer use
// ADR: 20260328-vision-pipeline-computer-use
type ScreenState struct {
	mu      sync.RWMutex
	History []ScreenSnapshot
	Cursor  Point
	Focus   string  // Currently focused element
	ScrollY int
	URL     string
}

// ScreenSnapshot represents a single screen state
type ScreenSnapshot struct {
	Timestamp time.Time
	Base64    string
	Hash      string  // SHA256 for change detection
	Width     int
	Height    int
}

// Point represents cursor/element position
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// NewScreenState creates a new screen state tracker
func NewScreenState() *ScreenState {
	return &ScreenState{
		History: make([]ScreenSnapshot, 0, 100), // Keep last 100 states
	}
}

// AddSnapshot adds a new screenshot to the history
func (s *ScreenState) AddSnapshot(base64Data string, width, height int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash := hashScreenshot(base64Data)
	snapshot := ScreenSnapshot{
		Timestamp: time.Now(),
		Base64:    base64Data,
		Hash:      hash,
		Width:     width,
		Height:    height,
	}

	s.History = append(s.History, snapshot)

	// Prune old snapshots (keep last 100)
	if len(s.History) > 100 {
		s.History = s.History[len(s.History)-100:]
	}
}

// HasChanged checks if the new screenshot is different from the last
func (s *ScreenState) HasChanged(newBase64 string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.History) == 0 {
		return true
	}

	lastHash := s.History[len(s.History)-1].Hash
	newHash := hashScreenshot(newBase64)
	return lastHash != newHash
}

// GetLastSnapshot returns the most recent screen snapshot
func (s *ScreenState) GetLastSnapshot() *ScreenSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.History) == 0 {
		return nil
	}
	return &s.History[len(s.History)-1]
}

// GetHistory returns the full screen history
func (s *ScreenState) GetHistory() []ScreenSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]ScreenSnapshot, len(s.History))
	copy(result, s.History)
	return result
}

// Prune removes old snapshots beyond maxHistory
func (s *ScreenState) Prune(maxHistory int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.History) > maxHistory {
		s.History = s.History[len(s.History)-maxHistory:]
	}
}

// SetCursor updates the cursor position
func (s *ScreenState) SetCursor(x, y int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Cursor = Point{X: x, Y: y}
}

// SetFocus updates the focused element
func (s *ScreenState) SetFocus(element string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Focus = element
}

// SetURL updates the current URL
func (s *ScreenState) SetURL(url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.URL = url
}

// GetStateJSON returns the current state as JSON
func (s *ScreenState) GetStateJSON() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type StateJSON struct {
		URL      string    `json:"url"`
		Cursor   Point     `json:"cursor"`
		Focus    string    `json:"focus"`
		ScrollY  int       `json:"scroll_y"`
		Screenshot *ScreenSnapshot `json:"last_screenshot,omitempty"`
	}

	state := StateJSON{
		URL:     s.URL,
		Cursor:  s.Cursor,
		Focus:   s.Focus,
		ScrollY: s.ScrollY,
	}
	if len(s.History) > 0 {
		state.Screenshot = &s.History[len(s.History)-1]
	}

	data, _ := json.Marshal(state)
	return string(data)
}

// Clear resets the screen state
func (s *ScreenState) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.History = make([]ScreenSnapshot, 0, 100)
	s.URL = ""
	s.Focus = ""
	s.Cursor = Point{}
	s.ScrollY = 0
}

// hashScreenshot computes SHA256 hash for screenshot comparison
func hashScreenshot(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
