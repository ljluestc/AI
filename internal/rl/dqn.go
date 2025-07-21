package rl

import (
	"fmt"
	"math/rand"

	"gonum.org/v1/gonum/mat"
)

// DQN represents a Deep Q-Network for reinforcement learning
type DQN struct {
	StateDim      int
	ActionDim     int
	QNetwork      *mat.Dense
	TargetNetwork *mat.Dense
	Memory        *BuddyMemory
	Optimizer     *AdamOptimizer
}

// ActionType represents an action that the agent can take
type ActionType struct {
	ID          int // Use int for ID consistently
	Name        string
	Description string
	Parameters  []float64
}

// NewDQN creates a new DQN with the given state and action dimensions
func NewDQN(stateDim, actionDim int, memory *BuddyMemory) *DQN {
	return &DQN{
		StateDim:      stateDim,
		ActionDim:     actionDim,
		QNetwork:      mat.NewDense(stateDim, actionDim, nil),
		TargetNetwork: mat.NewDense(stateDim, actionDim, nil),
		Memory:        memory,
		Optimizer:     NewAdamOptimizer(0.001, 0.9, 0.999),
	}
}

// SelectAction selects an action based on the current state
func (dqn *DQN) SelectAction(state []float64, epsilon float64) (ActionType, error) {
	if len(state) != dqn.StateDim {
		return ActionType{}, fmt.Errorf("invalid state dimension: expected %d, got %d", dqn.StateDim, len(state))
	}

	// Epsilon-greedy action selection
	if rand.Float64() < epsilon {
		// Explore: select a random action
		actionID := rand.Intn(dqn.ActionDim)
		return ActionType{
			ID:   actionID,
			Name: fmt.Sprintf("Action_%d", actionID),
		}, nil
	}

	// Exploit: select the best action according to the Q-network
	stateVector := mat.NewVecDense(dqn.StateDim, state)
	qValues := mat.NewVecDense(dqn.ActionDim, nil)
	qValues.MulVec(dqn.QNetwork.T(), stateVector)

	// Find the action with the highest Q-value
	bestActionID := 0
	bestValue := qValues.AtVec(0)
	for i := 1; i < dqn.ActionDim; i++ {
		if qValues.AtVec(i) > bestValue {
			bestValue = qValues.AtVec(i)
			bestActionID = i
		}
	}

	return ActionType{
		ID:   bestActionID,
		Name: fmt.Sprintf("Action_%d", bestActionID),
	}, nil
}

// Update updates the DQN based on a batch of experiences
func (dqn *DQN) Update(gamma float64) error {
	// Implementation would depend on how experiences are stored and processed
	// This is a simplified version that just updates the target network
	dqn.TargetNetwork.Copy(dqn.QNetwork)
	return nil
}
