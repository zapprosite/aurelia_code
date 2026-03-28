package audio

import (
	"errors"
	"io"
	"sync"
)

var (
	ErrBufferFull = errors.New("audio buffer full")
)

// StreamBuffer é um buffer circular thread-safe para chunks de áudio PCM.
// Ideal para alimentar o Voice Loop sem latência de disco.
type StreamBuffer struct {
	mu     sync.RWMutex
	data   []byte
	head   int
	tail   int
	size   int
	isFull bool
}

func NewStreamBuffer(capacity int) *StreamBuffer {
	return &StreamBuffer{
		data: make([]byte, capacity),
		size: capacity,
	}
}

// Write implementa io.Writer para o buffer circular.
func (b *StreamBuffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, b_byte := range p {
		if b.isFull {
			return n, ErrBufferFull
		}
		b.data[b.head] = b_byte
		b.head = (b.head + 1) % b.size
		n++
		if b.head == b.tail {
			b.isFull = true
		}
	}
	return n, nil
}

// Read implementa io.Reader para o buffer circular.
func (b *StreamBuffer) Read(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.head == b.tail && !b.isFull {
		return 0, io.EOF
	}

	for i := range p {
		if b.head == b.tail && !b.isFull {
			break
		}
		p[i] = b.data[b.tail]
		b.tail = (b.tail + 1) % b.size
		b.isFull = false
		n++
	}
	return n, nil
}

// Reset limpa o buffer.
func (b *StreamBuffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.head = 0
	b.tail = 0
	b.isFull = false
}
