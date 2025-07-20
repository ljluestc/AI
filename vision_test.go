package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

// TestVisionTransformerInitialization tests the VisionTransformer initialization
func TestVisionTransformerInitialization(t *testing.T) {
	// Test with depth enabled
	vt := &VisionTransformer{
		numFrames:     5,
		useDepth:      true,
		inputChannels: 4,
		weights:       mat.NewDense(768, 1000, nil),
	}
	assert.Equal(t, 5, vt.numFrames)
	assert.True(t, vt.useDepth)
	assert.Equal(t, 4, vt.inputChannels)
	
	// Test with depth disabled
	vt = &VisionTransformer{
		numFrames:     3,
		useDepth:      false,
		inputChannels: 3,
		weights:       mat.NewDense(768, 1000, nil),
	}
	assert.Equal(t, 3, vt.numFrames)
	assert.False(t, vt.useDepth)
	assert.Equal(t, 3, vt.inputChannels)
}

// TestVisionInputProcessing tests various input formats for the vision module
func TestVisionInputProcessing(t *testing.T) {
	vt := &VisionTransformer{
		numFrames:     3,
		useDepth:      true,
		inputChannels: 4,
		weights:       mat.NewDense(768, 1000, nil),
	}
	
	// Test with empty input
	emptyInput := VisionInput{
		RGBFrames:   [][][]float64{},
		DepthFrames: [][][]float64{},
	}
	state, err := vt.ProcessVisionInput(emptyInput)
	assert.NoError(t, err)
	assert.NotNil(t, state)
	assert.Len(t, state.Position, 3)
	
	// Test with single frame input
	singleFrameInput := VisionInput{
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
	state, err = vt.ProcessVisionInput(singleFrameInput)
	assert.NoError(t, err)
	assert.NotNil(t, state)
	assert.Len(t, state.Position, 3)
	
	// Test with multi-frame input
	multiFrameInput := VisionInput{
		RGBFrames: [][][]float64{
			{
				{1.0, 2.0, 3.0},
				{4.0, 5.0, 6.0},
			},
			{
				{7.0, 8.0, 9.0},
				{10.0, 11.0, 12.0},
			},
			{
				{13.0, 14.0, 15.0},
				{16.0, 17.0, 18.0},
			},
		},
		DepthFrames: [][][]float64{
			{
				{0.5, 0.6, 0.7},
				{0.8, 0.9, 1.0},
			},
			{
				{1.1, 1.2, 1.3},
				{1.4, 1.5, 1.6},
			},
			{
				{1.7, 1.8, 1.9},
				{2.0, 2.1, 2.2},
			},
		},
	}
	state, err = vt.ProcessVisionInput(multiFrameInput)
	assert.NoError(t, err)
	assert.NotNil(t, state)
	assert.Len(t, state.Position, 3)
}

// TestSegmentationMaskGeneration tests the segmentation mask generation
func TestSegmentationMaskGeneration(t *testing.T) {
	vt := &VisionTransformer{
		numFrames:     3,
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
	
	state, err := vt.ProcessVisionInput(input)
	assert.NoError(t, err)
	
	// Check mask dimensions
	assert.Equal(t, vt.numFrames, len(state.Mask))
	assert.Equal(t, 224, len(state.Mask[0]))
	assert.Equal(t, 224, len(state.Mask[0][0]))
	
	// Check mask values are in [0,1] range
	for i := 0; i < vt.numFrames; i++ {
		for j := 0; j < 224; j++ {
			for k := 0; k < 224; k++ {
				assert.GreaterOrEqual(t, state.Mask[i][j][k], 0.0)
				assert.LessOrEqual(t, state.Mask[i][j][k], 1.0)
			}
		}
	}
}

// TestDepthAwareProcessing tests the depth-aware processing
func TestDepthAwareProcessing(t *testing.T) {
	// Test with depth enabled
	vtWithDepth := &VisionTransformer{
		numFrames:     3,
		useDepth:      true,
		inputChannels: 4,
		weights:       mat.NewDense(768, 1000, nil),
	}
	
	// Test with depth disabled
	vtWithoutDepth := &VisionTransformer{
		numFrames:     3,
		useDepth:      false,
		inputChannels: 3,
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
	
	// Process with depth enabled
	stateWithDepth, err := vtWithDepth.ProcessVisionInput(input)
	assert.NoError(t, err)
	
	// Process with depth disabled
	stateWithoutDepth, err := vtWithoutDepth.ProcessVisionInput(input)
	assert.NoError(t, err)
	
	// Both should return valid results
	assert.NotNil(t, stateWithDepth)
	assert.NotNil(t, stateWithoutDepth)
	assert.Len(t, stateWithDepth.Position, 3)
	assert.Len(t, stateWithoutDepth.Position, 3)
}
