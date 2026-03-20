package agent

import (
"sync"
)

type EventKind string

const (
EventTaskCreated   EventKind = "task_created"
EventTaskAssigned  EventKind = "task_assigned"
EventHelpRequested EventKind = "help_requested"
EventHelpOffered   EventKind = "help_offered"
EventTaskCompleted EventKind = "task_completed"
)

type Event struct {
Kind    EventKind
Payload interface{}
TeamID  string
RunID   string
}

type Subscriber func(Event)

type EventBus struct {
mu          sync.RWMutex
subscribers map[EventKind][]Subscriber
}

func NewEventBus() *EventBus {
return &EventBus{
tKind][]Subscriber),
}
}

func (b *EventBus) Subscribe(kind EventKind, sub Subscriber) {
b.mu.Lock()
defer b.mu.Unlock()
b.subscribers[kind] = append(b.subscribers[kind], sub)
}

func (b *EventBus) Publish(event Event) {
b.mu.RLock()
defer b.mu.RUnlock()
for _, sub := range b.subscribers[event.Kind] {
t)
}
}
