package events

import "sync"

type Broadcaster struct {
	clients map[chan string]bool
	mu      sync.Mutex
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[chan string]bool),
	}
}

func (b *Broadcaster) Subscribe() chan string {
	ch := make(chan string, 1)
	b.mu.Lock()
	b.clients[ch] = true
	b.mu.Unlock()
	return ch
}

func (b *Broadcaster) Unsubscribe(ch chan string) {
	b.mu.Lock()
	delete(b.clients, ch)
	close(ch)
	b.mu.Unlock()
}

func (b *Broadcaster) Publish(msg string) {
	b.mu.Lock()
	for ch := range b.clients {
		select {
		case ch <- msg:
		default:
			// if channel is full, drop message
		}
	}
	b.mu.Unlock()
}
