package rl

// ReplayMemory stores experiences for reinforcement learning
// This is a facade that delegates to the memory package
type ReplayMemory struct {
	capacity    int
	experiences []Experience
	counter     int
}

// Experience represents a single learning experience
type Experience struct {
	State     []float64
	Action    int
	Reward    float64
	NextState []float64
	IsDone    bool
}

// NewReplayMemory creates a new ReplayMemory with the given capacity
func NewReplayMemory(capacity int) *ReplayMemory {
	return &ReplayMemory{
		capacity:    capacity,
		experiences: make([]Experience, 0, capacity),
		counter:     0,
	}
}

// Add adds an experience to memory
func (m *ReplayMemory) Add(state []float64, action int, reward float64, nextState []float64, isDone bool) {
	exp := Experience{
		State:     state,
		Action:    action,
		Reward:    reward,
		NextState: nextState,
		IsDone:    isDone,
	}

	if len(m.experiences) < m.capacity {
		m.experiences = append(m.experiences, exp)
	} else {
		// Circular buffer replacement
		idx := m.counter % m.capacity
		m.experiences[idx] = exp
	}
	m.counter++
}

// Sample samples a batch of experiences from memory
func (m *ReplayMemory) Sample(batchSize int) []Experience {
	if len(m.experiences) == 0 {
		return []Experience{}
	}

	if batchSize > len(m.experiences) {
		batchSize = len(m.experiences)
	}

	result := make([]Experience, batchSize)
	for i := 0; i < batchSize; i++ {
		idx := (m.counter - i - 1) % len(m.experiences)
		if idx < 0 {
			idx += len(m.experiences)
		}
		result[i] = m.experiences[idx]
	}

	return result
}
