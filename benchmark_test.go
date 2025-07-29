package main

import (
	"testing"

	"gonum.org/v1/gonum/mat"

	"github.com/teathis/codeanalyzer/internal/agent"
	"github.com/teathis/codeanalyzer/internal/config"
	"github.com/teathis/codeanalyzer/internal/knowledge"
	"github.com/teathis/codeanalyzer/internal/rl"
	"github.com/teathis/codeanalyzer/internal/vision"
)

// BenchmarkVisionTransformer benchmarks the vision processing performance
func BenchmarkVisionTransformer(b *testing.B) {
	vt := vision.NewVisionTransformer(5, true)
	input := vision.VisionInput{
		RGBFrames:   [][][][]float64{{{{0.1, 0.2, 0.3}}}},
		DepthFrames: [][][][]float64{{{{0.5}}}},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := vt.Process(input)
		if err != nil {
			b.Fatalf("Vision processing failed: %v", err)
		}
	}
}

// BenchmarkDQN benchmarks the DQN action selection performance
func BenchmarkDQN(b *testing.B) {
	memory := &rl.BuddyMemory{} // Removed Allocated field
	dqn := rl.NewDQN(100, 10, memory)
	state := make([]float64, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := dqn.SelectAction(state, 0.1)
		if err != nil {
			b.Fatalf("Action selection failed: %v", err)
		}
	}
}

// BenchmarkBuddyMemoryAllocation benchmarks the buddy memory allocation performance
func BenchmarkBuddyMemory(b *testing.B) {
	bm := &rl.BuddyMemory{} // Removed Allocated field
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		addr, err := bm.Allocate(1024)
		if err != nil {
			b.Fatalf("Allocation failed: %v", err)
		}
		bm.Free(addr)
	}
}

// BenchmarkAdamOptimizer benchmarks the Adam optimizer performance
func BenchmarkAdamOptimizer(b *testing.B) {
	optimizer := rl.NewAdamOptimizer(0.001, 0.9, 0.999)

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

// BenchmarkDebugAgent benchmarks the full debugging pipeline
func BenchmarkDebugAgent(b *testing.B) {
	cfg := config.Config{
		Port:          "8080",
		Neo4jURI:      "bolt://localhost:7687",
		Neo4jUser:     "neo4j",
		Neo4jPassword: "password",
	}
	ag, err := agent.NewDebugAgent(cfg)
	if err != nil {
		b.Fatalf("Failed to initialize DebugAgent: %v", err)
	}

	var visionInput = vision.VisionInput{
		RGBFrames: [][][][]float64{
			{
				{
					{1.0, 2.0, 3.0},
				},
			},
		},
		DepthFrames: [][][][]float64{
			{
				{
					{0.5, 0.6, 0.7},
					{0.8, 0.9, 1.0},
				},
			},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Step 1: Vision
		state, err := ag.VisionModel.Process(visionInput)
		if err != nil {
			b.Fatal(err)
		}

		// Step 2: Diagnosis (mock to avoid Neo4j dependency)
		diagnosisResult := knowledge.DiagnosisResult{
			Issue:        "cursor_lag",
			Confidence:   0.8,
			RelatedNodes: []string{"hardware:mouse", "software:driver"},
		}

		// Step 3: Planning
		action, err := ag.RLModel.SelectAction(state.Position[:], 0.1)
		if err != nil {
			b.Fatal(err)
		}

		// Step 4: Execution (mock)
		executionResult := agent.ExecutionResult{
			Success: true,
			Message: "Action executed successfully",
		}

		_ = diagnosisResult
		_ = action
		_ = executionResult
	}
}

// BenchmarkDiagnosisResult benchmarks the DiagnosisResult struct
func BenchmarkDiagnosisResult(b *testing.B) {
	result := knowledge.DiagnosisResult{
		Issue:        "KeyError",
		Confidence:   0.95,
		RelatedNodes: []string{"DataFrame", "Index"},
	}
	_ = result
}
