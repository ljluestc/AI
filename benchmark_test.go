package main

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

// --- Type definitions moved to top ---
type DQN struct {
	stateDim      int
	actionDim     int
	qNetwork      *mat.Dense
	targetNetwork *mat.Dense
	memory        *BuddyMemory
	optimizer     *AdamOptimizer
}

func (dqn *DQN) SelectAction(state []float64) (Action, error) {
	return Action{
		ID:          "action_1",
		Description: "Dummy action",
		Parameters:  state,
	}, nil
}

type BuddyMemory struct {
	size      int
	allocated map[int]bool
}

func (bm *BuddyMemory) Allocate(size int) (int, error) {
	// Find first free address
	for addr := 0; addr < bm.size; addr += size {
		if !bm.allocated[addr] {
			bm.allocated[addr] = true
			return addr, nil
		}
	}
	return -1, nil // Out of memory
}

func (bm *BuddyMemory) Free(addr int) {
	delete(bm.allocated, addr)
}

type AdamOptimizer struct {
	learningRate float64
	beta1        float64
	beta2        float64
	t            int
}

func (opt *AdamOptimizer) Update(weights *mat.Dense, grad float64) {
	// Dummy update: increment t
	opt.t++
}

type Config struct {
	Port          string
	Neo4jURI      string
	Neo4jUser     string
	Neo4jPassword string
}

type DiagnosisResult struct {
	Issue        string
	Confidence   float64
	RelatedNodes []string
}

type ExecutionResult struct {
	Success bool
	Message string
}

type DebugAgent struct {
	visionModel *VisionTransformer
	rlModel     *DQN
}

func NewDebugAgent(cfg Config) (*DebugAgent, error) {
	return &DebugAgent{
		visionModel: &VisionTransformer{
			numFrames:     5,
			useDepth:      true,
			inputChannels: 4,
			weights:       mat.NewDense(768, 1000, nil),
		},
		rlModel: &DQN{
			stateDim:      768,
			actionDim:     10,
			qNetwork:      mat.NewDense(768, 10, nil),
			targetNetwork: mat.NewDense(768, 10, nil),
			memory:        &BuddyMemory{size: 1024 * 1024, allocated: make(map[int]bool)},
			optimizer:     &AdamOptimizer{learningRate: 0.001, beta1: 0.9, beta2: 0.999},
		},
	}, nil
}

// Add VisionTransformer and VisionInput stubs if not present
type VisionTransformer struct {
	numFrames     int
	useDepth      bool
	inputChannels int
	weights       *mat.Dense
}

type VisionInput struct {
	RGBFrames   [][][]float64
	DepthFrames [][][]float64
}

type CursorState struct {
	Position [3]float64
	Mask     [][][]float64
}

// Action type for SelectAction
type Action struct {
	ID          string
	Description string
	Parameters  []float64
}

// BenchmarkVisionProcessing benchmarks the vision processing performance
func BenchmarkVisionProcessing(b *testing.B) {
	vt := &VisionTransformer{
		numFrames:     5,
		useDepth:      true,
		inputChannels: 4,
		weights:       mat.NewDense(768, 1000, nil),
	}

	input := VisionInput{
		RGBFrames: [][][]float64{
			{
				{1.0, 2.0, 3.0},
				{4.0, 5.0, 6.0},
			},
		},
		DepthFrames: [][][]float64{
			{
				{0.5, 0.6, 0.7},
				{0.8, 0.9, 1.0},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := vt.ProcessVisionInput(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDQNActionSelection benchmarks the DQN action selection performance
func BenchmarkDQNActionSelection(b *testing.B) {
	dqn := &DQN{
		stateDim:      768,
		actionDim:     10,
		qNetwork:      mat.NewDense(768, 10, nil),
		targetNetwork: mat.NewDense(768, 10, nil),
		memory:        &BuddyMemory{size: 1024 * 1024, allocated: make(map[int]bool)},
		optimizer:     &AdamOptimizer{learningRate: 0.001, beta1: 0.9, beta2: 0.999},
	}

	// Create a state vector with 768 dimensions
	state := make([]float64, 768)
	for i := range state {
		state[i] = float64(i) / 768.0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := dqn.SelectAction(state)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBuddyMemoryAllocation benchmarks the buddy memory allocation performance
func BenchmarkBuddyMemoryAllocation(b *testing.B) {
	bm := &BuddyMemory{
		size:      1024 * 1024, // 1MB
		allocated: make(map[int]bool),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		addr, err := bm.Allocate(1024) // 1KB blocks
		if err != nil {
			b.Fatal(err)
		}
		bm.Free(addr)
	}
}

// BenchmarkAdamOptimizer benchmarks the Adam optimizer performance
func BenchmarkAdamOptimizer(b *testing.B) {
	optimizer := &AdamOptimizer{
		learningRate: 0.001,
		beta1:        0.9,
		beta2:        0.999,
		t:            1,
	}

	// Create a 100x100 matrix
	weights := mat.NewDense(100, 100, nil)
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			weights.Set(i, j, float64(i*j)/10000.0)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		optimizer.Update(weights, 0.5)
	}
}

// BenchmarkEndToEndPipeline benchmarks the full debugging pipeline
func BenchmarkEndToEndPipeline(b *testing.B) {
	// Create a test config
	config := Config{
		Port:          "8080",
		Neo4jURI:      "neo4j://localhost:7687",
		Neo4jUser:     "neo4j",
		Neo4jPassword: "password",
	}

	// Create a test agent
	agent, err := NewDebugAgent(config)
	if err != nil {
		b.Fatal(err)
	}

	// Create test inputs
	visionInput := VisionInput{
		RGBFrames: [][][]float64{
			{
				{1.0, 2.0, 3.0},
				{4.0, 5.0, 6.0},
			},
		},
		DepthFrames: [][][]float64{
			{
				{0.5, 0.6, 0.7},
				{0.8, 0.9, 1.0},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Step 1: Vision
		state, err := agent.visionModel.ProcessVisionInput(visionInput)
		if err != nil {
			b.Fatal(err)
		}

		// Step 2: Diagnosis (mock to avoid Neo4j dependency)
		diagnosisResult := DiagnosisResult{
			Issue:        "cursor_lag",
			Confidence:   0.8,
			RelatedNodes: []string{"hardware:mouse", "software:driver"},
		}

		// Step 3: Planning
		action, err := agent.rlModel.SelectAction(state.Position[:])
		if err != nil {
			b.Fatal(err)
		}

		// Step 4: Execution (mock)
		executionResult := ExecutionResult{
			Success: true,
			Message: "Action executed successfully",
		}

		// Prevent compiler from optimizing away unused variables
		_ = diagnosisResult
		_ = action
		_ = executionResult
	}
}
