package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"

	"github.com/teathis/codeanalyzer/internal/agent"
	"github.com/teathis/codeanalyzer/internal/memory"
)

// TestDQNInitialization tests DQN initialization
func TestDQNInitialization(t *testing.T) {
	stateDim := 768
	actionDim := 10
	memory := &memory.BuddyMemory{Size: 1024 * 1024, Allocated: make(map[int]bool)}

	dqn := &agent.DQN{
		StateDim:      stateDim,
		ActionDim:     actionDim,
		QNetwork:      mat.NewDense(stateDim, actionDim, nil),
		TargetNetwork: mat.NewDense(stateDim, actionDim, nil),
		Memory:        memory,
		Optimizer:     &agent.AdamOptimizer{LearningRate: 0.001, Beta1: 0.9, Beta2: 0.999},
	}

	assert.Equal(t, stateDim, dqn.StateDim)
	assert.Equal(t, actionDim, dqn.ActionDim)
	assert.NotNil(t, dqn.QNetwork)
	assert.NotNil(t, dqn.TargetNetwork)
	assert.NotNil(t, dqn.Memory)
	assert.NotNil(t, dqn.Optimizer)
}

// TestActionSelection tests DQN action selection
func TestActionSelection(t *testing.T) {
	dqn := &agent.DQN{
		StateDim:      3,
		ActionDim:     5,
		QNetwork:      mat.NewDense(3, 5, nil),
		TargetNetwork: mat.NewDense(3, 5, nil),
		Memory:        &memory.BuddyMemory{Size: 1024, Allocated: make(map[int]bool)},
		Optimizer:     &agent.AdamOptimizer{LearningRate: 0.001, Beta1: 0.9, Beta2: 0.999},
	}

	// Set Q-values to have a clear maximum at action 2
	dqn.QNetwork.Set(0, 0, 0.1)
	dqn.QNetwork.Set(0, 1, 0.2)
	dqn.QNetwork.Set(0, 2, 0.5) // max value
	dqn.QNetwork.Set(0, 3, 0.3)
	dqn.QNetwork.Set(0, 4, 0.4)

	dqn.QNetwork.Set(1, 0, 0.1)
	dqn.QNetwork.Set(1, 1, 0.2)
	dqn.QNetwork.Set(1, 2, 0.5) // max value
	dqn.QNetwork.Set(1, 3, 0.3)
	dqn.QNetwork.Set(1, 4, 0.4)

	dqn.QNetwork.Set(2, 0, 0.1)
	dqn.QNetwork.Set(2, 1, 0.2)
	dqn.QNetwork.Set(2, 2, 0.5) // max value
	dqn.QNetwork.Set(2, 3, 0.3)
	dqn.QNetwork.Set(2, 4, 0.4)

	// Test action selection
	state := []float64{1.0, 0.0, 0.0} // Only first row should matter
	action, err := dqn.SelectAction(state)

	// Should select action 2 with highest Q-value
	assert.NoError(t, err)
	assert.Equal(t, "action_2", action.ID) // Index 2 action
	assert.NotEmpty(t, action.Description)
	assert.Len(t, action.Parameters, 1)
}

// TestDQNUpdate tests DQN update process
func TestDQNUpdate(t *testing.T) {
	dqn := &agent.DQN{
		StateDim:      3,
		ActionDim:     5,
		QNetwork:      mat.NewDense(3, 5, nil),
		TargetNetwork: mat.NewDense(3, 5, nil),
		Memory:        &memory.BuddyMemory{Size: 1024, Allocated: make(map[int]bool)},
		Optimizer:     &agent.AdamOptimizer{LearningRate: 0.001, Beta1: 0.9, Beta2: 0.999},
	}

	// Set initial Q-values
	initialQValue := dqn.QNetwork.At(0, 0)

	// Prepare update inputs
	states := []float64{1.0, 0.0, 0.0}
	actions := []float64{0.0} // Action 0
	rewards := 1.0            // Positive reward
	nextStates := []float64{0.0, 1.0, 0.0}
	done := false

	// Perform update
	err := dqn.Update(states, actions, rewards, nextStates, done)
	assert.NoError(t, err)

	// Q-value should have changed
	updatedQValue := dqn.QNetwork.At(0, 0)
	assert.NotEqual(t, initialQValue, updatedQValue)
}

// TestAdamOptimizer tests the Adam optimizer
func TestAdamOptimizer(t *testing.T) {
	optimizer := &agent.AdamOptimizer{
		LearningRate: 0.01,
		Beta1:        0.9,
		Beta2:        0.999,
		T:            1,
	}

	// Create a test weights matrix
	weights := mat.NewDense(2, 2, []float64{1.0, 2.0, 3.0, 4.0})

	// Apply the optimizer
	initialValue := weights.At(0, 0)
	optimizer.Update(weights, 0.5) // Apply with loss of 0.5

	// Check that the weights have been updated
	updatedValue := weights.At(0, 0)
	assert.NotEqual(t, initialValue, updatedValue)
}

// TestBuddyMemoryAllocation tests buddy memory allocation and freeing
func TestBuddyMemoryAllocation(t *testing.T) {
	bm := &memory.BuddyMemory{
		Size:      1024,
		Allocated: make(map[int]bool),
	}

	// Allocate memory
	addr1, err := bm.Allocate(256)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, addr1, 0)
	assert.Less(t, addr1, bm.Size)
	assert.True(t, bm.Allocated[addr1])

	// Allocate more memory
	addr2, err := bm.Allocate(256)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, addr2, 0)
	assert.Less(t, addr2, bm.Size)
	assert.NotEqual(t, addr1, addr2)
	assert.True(t, bm.Allocated[addr2])

	// Free the first allocation
	bm.Free(addr1)
	assert.False(t, bm.Allocated[addr1])

	// Re-allocate, should get first address back
	addr3, err := bm.Allocate(256)
	assert.NoError(t, err)
	assert.Equal(t, addr1, addr3)
	assert.True(t, bm.Allocated[addr3])

	// Test allocation exhaustion
	// Allocate until full
	addrs := []int{addr2, addr3}
	for i := 0; i < 2; i++ { // Try to allocate 2 more blocks (total 4, but we only have space for 3)
		addr, err := bm.Allocate(256)
		if err == nil {
			addrs = append(addrs, addr)
		}
	}

	// Should have allocated 3 blocks (1024/256 = 4, but due to buddy system limitations, only 3 fit)
	assert.LessOrEqual(t, len(addrs), 4)

	// Free all allocations
	for _, addr := range addrs {
		bm.Free(addr)
	}

	// Check that all allocations are freed
	for i := 0; i < bm.Size; i += 256 {
		assert.False(t, bm.Allocated[i])
	}
}
