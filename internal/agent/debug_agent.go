package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/teathis/codeanalyzer/internal/config"
	"github.com/teathis/codeanalyzer/internal/knowledge"
	"github.com/teathis/codeanalyzer/internal/rl"
	"github.com/teathis/codeanalyzer/internal/vision"
)

// DebugAgent is the main component responsible for coordinating the debugging process
type DebugAgent struct {
	VisionModel    *vision.VisionTransformer
	RLModel        *rl.DQN
	KnowledgeGraph *knowledge.KnowledgeGraph
	Logger         *logrus.Logger
	Clients        map[*websocket.Conn]bool
	ClientsMu      sync.Mutex
}

// ExecutionResult represents the result of executing an action
type ExecutionResult struct {
	Success bool
	Message string
}

// NewDebugAgent creates a new DebugAgent with the given configuration
func NewDebugAgent(cfg config.Config) (*DebugAgent, error) {
	kg, err := knowledge.NewKnowledgeGraph(cfg.Neo4jURI, cfg.Neo4jUser, cfg.Neo4jPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to create knowledge graph: %w", err)
	}

	return &DebugAgent{
		VisionModel:    vision.NewVisionTransformer(5, true),
		RLModel:        rl.NewDQN(100, 10, rl.NewBuddyMemory(1024*1024)),
		KnowledgeGraph: kg,
		Logger:         logrus.New(),
		Clients:        make(map[*websocket.Conn]bool),
	}, nil
}

// HandleVision processes the vision input and returns the cursor state
func (a *DebugAgent) HandleVision(input vision.VisionInput) (*vision.CursorState, error) {
	return a.VisionModel.Process(input)
}

// HandleDiagnosis processes the cursor state and returns the diagnosis result
func (a *DebugAgent) HandleDiagnosis(state *vision.CursorState) (*knowledge.DiagnosisResult, error) {
	return a.KnowledgeGraph.DiagnoseCursorIssue(state)
}

// HandlePlanning processes the diagnosis result and returns the action to take
func (a *DebugAgent) HandlePlanning(state *vision.CursorState) (rl.Action, error) {
	return a.RLModel.SelectAction(state.Position[:], 0.1)
}

// HandleExecution executes the action and returns the result
func (a *DebugAgent) HandleExecution(action rl.Action) (*ExecutionResult, error) {
	// Implementation details would depend on what actions need to be executed
	return &ExecutionResult{
		Success: true,
		Message: fmt.Sprintf("Executed action %s", action.Name),
	}, nil
}

// ProcessDebugSession processes a complete debug session from vision to execution
func (a *DebugAgent) ProcessDebugSession(ctx context.Context, input vision.VisionInput) (*ExecutionResult, error) {
	// Step 1: Vision
	state, err := a.HandleVision(input)
	if err != nil {
		return nil, fmt.Errorf("vision processing failed: %w", err)
	}

	// Step 2: Diagnosis
	diagnosis, err := a.HandleDiagnosis(state)
	if err != nil {
		return nil, fmt.Errorf("diagnosis failed: %w", err)
	}

	// Step 3: Planning
	action, err := a.HandlePlanning(state)
	if err != nil {
		return nil, fmt.Errorf("planning failed: %w", err)
	}

	// Step 4: Execution
	result, err := a.HandleExecution(action)
	if err != nil {
		return nil, fmt.Errorf("execution failed: %w", err)
	}

	// Log the results
	a.Logger.WithFields(logrus.Fields{
		"cursor_state": state,
		"diagnosis":    diagnosis,
		"action":       action,
		"result":       result,
	}).Info("Debug session completed")

	return result, nil
}
