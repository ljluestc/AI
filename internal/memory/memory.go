// Package memory provides memory management for the RL agent
package memory

import (
	"errors"
	"math/rand"
)

// Experience represents a single learning experience
type Experience struct {
	State     []float64
	Action    int
	Reward    float64
	NextState []float64
	IsDone    bool
}

// ReplayBuffer is a simple experience replay buffer
type ReplayBuffer struct {
	capacity    int
	experiences []Experience
	position    int
	count       int
}

// NewReplayBuffer creates a new replay buffer with the given capacity
func NewReplayBuffer(capacity int) *ReplayBuffer {
	return &ReplayBuffer{
		capacity:    capacity,
		experiences: make([]Experience, capacity),
		position:    0,
		count:       0,
	}
}

// Add adds a new experience to the buffer
func (rb *ReplayBuffer) Add(state []float64, action int, reward float64, nextState []float64, isDone bool) {
	experience := Experience{
		State:     state,
		Action:    action,
		Reward:    reward,
		NextState: nextState,
		IsDone:    isDone,
	}

	rb.experiences[rb.position] = experience
	rb.position = (rb.position + 1) % rb.capacity
	if rb.count < rb.capacity {
		rb.count++
	}
}

// Sample randomly samples a batch of experiences from the buffer
func (rb *ReplayBuffer) Sample(batchSize int) []Experience {
	if rb.count == 0 {
		return []Experience{}
	}

	if batchSize > rb.count {
		batchSize = rb.count
	}

	// Generate random indices
	indices := rand.Perm(rb.count)[:batchSize]

	// Sample experiences
	batch := make([]Experience, batchSize)
	for i, idx := range indices {
		batch[i] = rb.experiences[idx]
	}

	return batch
}

// Size returns the current number of experiences in the buffer
func (rb *ReplayBuffer) Size() int {
	return rb.count
}

// Clear clears the buffer
func (rb *ReplayBuffer) Clear() {
	rb.count = 0
	rb.position = 0
}

// BuddyMemory stores and retrieves experiences for reinforcement learning
type BuddyMemory struct {
	capacity    int
	experiences []Experience
	counter     int
}

// NewBuddyMemory creates a new memory buffer with the specified capacity
func NewBuddyMemory(capacity int) *BuddyMemory {
	return &BuddyMemory{
		capacity:    capacity,
		experiences: make([]Experience, 0, capacity),
		counter:     0,
	}
}

// Add adds a new experience to memory
func (m *BuddyMemory) Add(exp Experience) error {
	if len(exp.State) == 0 || len(exp.NextState) == 0 {
		return errors.New("state vectors cannot be empty")
	}

	if len(m.experiences) < m.capacity {
		m.experiences = append(m.experiences, exp)
	} else {
		// Circular buffer replacement
		idx := m.counter % m.capacity
		m.experiences[idx] = exp
	}
	m.counter++
	return nil
}

// Sample randomly samples a batch of experiences from memory
func (m *BuddyMemory) Sample(batchSize int) ([]Experience, error) {
	if len(m.experiences) == 0 {
		return nil, errors.New("memory is empty")
	}

	if batchSize > len(m.experiences) {
		batchSize = len(m.experiences)
	}

	batch := make([]Experience, batchSize)
	indices := rand.Perm(len(m.experiences))[:batchSize]

	for i, idx := range indices {
		batch[i] = m.experiences[idx]
	}

	return batch, nil
}

// Size returns the current number of experiences in memory
func (m *BuddyMemory) Size() int {
	return len(m.experiences)
}

// Clear empties the memory
func (m *BuddyMemory) Clear() {
	m.experiences = make([]Experience, 0, m.capacity)
	m.counter = 0
}
