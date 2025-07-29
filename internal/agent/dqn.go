// Package agent provides agent implementations for the codeanalyzer
package agent

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"

	"github.com/teathis/codeanalyzer/internal/memory"
	"github.com/teathis/codeanalyzer/internal/rl"
	"github.com/teathis/codeanalyzer/pkg/optimizer"
)

// DQN represents a Deep Q-Network for reinforcement learning
type DQN struct {
	StateDim       int
	ActionDim      int
	QNetwork       mat.Matrix
	TargetQNetwork mat.Matrix
	Memory         *memory.BuddyMemory
	Optimizer      optimizer.Optimizer
}

// SelectAction selects an action based on the current state
func (dqn *DQN) SelectAction(state []float64) (rl.Action, error) {
	if len(state) != dqn.StateDim {
		return rl.Action{}, fmt.Errorf("invalid state dimension: expected %d, got %d", dqn.StateDim, len(state))
	}

	// Get Q-values for the state
	qNet := dqn.QNetwork.(*mat.Dense)

	// Find the action with the highest Q-value
	bestActionID := 0
	bestValue := qNet.At(0, 0)
	for i := 1; i < dqn.ActionDim; i++ {
		if qNet.At(0, i) > bestValue {
			bestValue = qNet.At(0, i)
			bestActionID = i
		}
	}

	// Return the selected action
	return rl.Action{
		ID:          fmt.Sprintf("action_%d", bestActionID),
		Description: fmt.Sprintf("Action %d", bestActionID),
		Parameters:  []float64{bestValue},
	}, nil
}

// Update updates the DQN based on the given experience
func (dqn *DQN) Update(states, actions []float64, rewards float64, nextStates []float64, done bool) error {
	if len(states) != dqn.StateDim {
		return fmt.Errorf("invalid state dimension: expected %d, got %d", dqn.StateDim, len(states))
	}
	if len(nextStates) != dqn.StateDim {
		return fmt.Errorf("invalid next state dimension: expected %d, got %d", dqn.StateDim, len(nextStates))
	}
	if len(actions) != 1 {
		return fmt.Errorf("invalid action dimension: expected 1, got %d", len(actions))
	}

	// Convert action to int
	actionID := int(actions[0])
	if actionID < 0 || actionID >= dqn.ActionDim {
		return fmt.Errorf("invalid action ID: %d", actionID)
	}

	// Get current Q-value for the state-action pair
	qNet := dqn.QNetwork.(*mat.Dense)
	currentQ := qNet.At(0, actionID)

	// Get max Q-value for the next state
	targetNet := dqn.TargetQNetwork.(*mat.Dense)
	maxNextQ := targetNet.At(0, 0)
	for i := 1; i < dqn.ActionDim; i++ {
		if targetNet.At(0, i) > maxNextQ {
			maxNextQ = targetNet.At(0, i)
		}
	}

	// Calculate target Q-value
	gamma := 0.99 // discount factor
	targetQ := rewards
	if !done {
		targetQ += gamma * maxNextQ
	}

	// Calculate loss (TD error)
	loss := math.Abs(targetQ - currentQ)

	// Update Q-network
	qNet.Set(0, actionID, targetQ)

	// Apply optimizer
	dqn.Optimizer.Update(qNet, loss)

	return nil
}
