package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"

	"github.com/teathis/codeanalyzer/internal/agent"
	"github.com/teathis/codeanalyzer/internal/memory"
	"github.com/teathis/codeanalyzer/pkg/optimizer"
)

// TestDQNInitialization tests DQN initialization
func TestDQNInitialization(t *testing.T) {
	stateDim := 768
	actionDim := 10
	memory := &memory.BuddyMemory{} // Initialize using methods instead of direct field access
	memory.Clear()                  // Clear will initialize the experiences slice

	dqn := &agent.DQN{
		StateDim:       stateDim,
		ActionDim:      actionDim,
		QNetwork:       mat.NewDense(stateDim, actionDim, nil),
		TargetQNetwork: mat.NewDense(stateDim, actionDim, nil),
		Memory:         memory,
		Optimizer:      optimizer.NewAdam(0.001, 0.9, 0.999),
	}

	assert.Equal(t, stateDim, dqn.StateDim)
	assert.Equal(t, actionDim, dqn.ActionDim)
	assert.NotNil(t, dqn.QNetwork)
	assert.NotNil(t, dqn.TargetQNetwork)
	assert.NotNil(t, dqn.Memory)
	assert.NotNil(t, dqn.Optimizer)
}

// TestActionSelection tests DQN action selection
func TestActionSelection(t *testing.T) {
	dqn := &agent.DQN{
		StateDim:       3,
		ActionDim:      5,
		QNetwork:       mat.NewDense(3, 5, nil),
		TargetQNetwork: mat.NewDense(3, 5, nil),
		Memory:         &memory.BuddyMemory{},
		Optimizer:      optimizer.NewAdam(0.001, 0.9, 0.999),
	}

	// Set Q-values to have a clear maximum at action 2
	qNet := dqn.QNetwork.(*mat.Dense)
	qNet.Set(0, 0, 0.1)
	qNet.Set(0, 1, 0.2)
	qNet.Set(0, 2, 0.5) // max value
	qNet.Set(0, 3, 0.3)
	qNet.Set(0, 4, 0.4)

	qNet.Set(1, 0, 0.1)
	qNet.Set(1, 1, 0.2)
	qNet.Set(1, 2, 0.5) // max value
	qNet.Set(1, 3, 0.3)
	qNet.Set(1, 4, 0.4)

	qNet.Set(2, 0, 0.1)
	qNet.Set(2, 1, 0.2)
	qNet.Set(2, 2, 0.5) // max value
	qNet.Set(2, 3, 0.3)
	qNet.Set(2, 4, 0.4)

	// Test action selection
	state := []float64{1.0, 0.0, 0.0} // Only first row should matter
	action, err := dqn.SelectAction(state)

	// Should select action 2 with highest Q-value
	assert.NoError(t, err)
	assert.Equal(t, "action_2", action.ID) // Index 2 action
	assert.NotEmpty(t, action.Desc)
	assert.Len(t, action.Parameters, 1)
}

// TestDQNUpdate tests DQN update process
func TestDQNUpdate(t *testing.T) {
	dqn := &agent.DQN{
		StateDim:       3,
		ActionDim:      5,
		QNetwork:       mat.NewDense(3, 5, nil),
		TargetQNetwork: mat.NewDense(3, 5, nil),
		Memory:         &memory.BuddyMemory{experiences: make([]memory.Experience, 0)},
		Optimizer:      optimizer.NewAdam(0.001, 0.9, 0.999),
	}

	// Set initial Q-values
	initialQValue := dqn.QNetwork.(*mat.Dense).At(0, 0)

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
	updatedQValue := dqn.QNetwork.(*mat.Dense).At(0, 0)
	assert.NotEqual(t, initialQValue, updatedQValue)
}

// TestAdamOptimizer tests the Adam optimizer
func TestAdamOptimizer(t *testing.T) {
	optimizer := optimizer.NewAdam(0.01, 0.9, 0.999)
	optimizer.T = 1

	// Create a test weights matrix
	weights := mat.NewDense(2, 2, []float64{1.0, 2.0, 3.0, 4.0})

	// Apply the optimizer
	initialValue := weights.At(0, 0)
	optimizer.Update(weights, 0.5) // Apply with loss of 0.5

	// Check that the weights have been updated
	updatedValue := weights.At(0, 0)
	assert.NotEqual(t, initialValue, updatedValue)
}

// TestBuddyMemory tests buddy memory operations
func TestBuddyMemory(t *testing.T) {
	bm := &memory.BuddyMemory{}
	// Set capacity if there's a method to do so, otherwise it might use a default value

	// Add an experience
	exp1 := memory.Experience{
		State:     []float64{1.0, 2.0},
		Action:    1,
		Reward:    0.5,
		NextState: []float64{2.0, 3.0},
		IsDone:    false,
	}

	err := bm.Add(exp1)
	assert.NoError(t, err)
	assert.Equal(t, 1, bm.Size())

	// Add another experience
	exp2 := memory.Experience{
		State:     []float64{2.0, 3.0},
		Action:    2,
		Reward:    1.0,
		NextState: []float64{3.0, 4.0},
		IsDone:    false,
	}

	err = bm.Add(exp2)
	assert.NoError(t, err)
	assert.Equal(t, 2, bm.Size())

	// Sample experiences
	batch, err := bm.Sample(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(batch))

	// Clear memory
	bm.Clear()
	assert.Equal(t, 0, bm.Size())
}
