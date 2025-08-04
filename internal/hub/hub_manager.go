package hub

import (
	"log"
	"sync"
	"z-chat/internal/domain/repository"
)

// Manager manages multiple chat hubs for different rooms.
type Manager struct {
	hubs map[string]*Hub
	mu   sync.RWMutex
	repo repository.MessageRepository
}

// NewManager creates a new Manager instance.
func NewManager(repo repository.MessageRepository) *Manager {
	return &Manager{
		hubs: make(map[string]*Hub),
		repo: repo,
	}
}

// GetOrCreateHub returns an existing hub for the given room ID or creates a new one if it doesn't exist.
func (m *Manager) GetOrCreateHub(roomID string) *Hub {
	m.mu.RLock()
	if hub, exists := m.hubs[roomID]; exists {
		m.mu.RUnlock()
		return hub
	}
	m.mu.RUnlock()

	// If we reach here, we need to create a new hub
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double check if another goroutine created the hub while we were waiting for the lock
	if hub, exists := m.hubs[roomID]; exists {
		return hub
	}

	// Create a new hub
	hub := NewHub(m.repo, roomID)
	m.hubs[roomID] = hub

	log.Printf("Creating new hub for roomID: %s", roomID)

	// Start the hub
	go hub.Run()

	return hub
}
