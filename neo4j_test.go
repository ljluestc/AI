package main

import (
	"context"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestKnowledgeGraphDiagnosis tests the diagnosis functionality
func TestKnowledgeGraphDiagnosis(t *testing.T) {
	// Create mock Neo4j driver and session
	mockDriver := new(MockNeo4jDriver)
	mockSession := new(MockNeo4jSession)

	// Create mock records and result
	mockRecords := []neo4j.Record{
		NewMockNeo4jRecord([]interface{}{"hardware:mouse", 0.8}, []string{"m.name", "r.confidence"}),
		NewMockNeo4jRecord([]interface{}{"software:driver", 0.7}, []string{"m.name", "r.confidence"}),
		NewMockNeo4jRecord([]interface{}{"system:memory", 0.6}, []string{"m.name", "r.confidence"}),
	}
	mockResult := NewMockNeo4jResult(mockRecords)

	// Set up expectations
	mockDriver.On("NewSession", mock.Anything).Return(mockSession)
	mockSession.On("Run", mock.Anything, mock.Anything).Return(mockResult, nil)
	mockSession.On("Close").Return(nil)
	mockResult.On("Err").Return(nil)
	mockResult.On("Consume").Return(nil, nil)

	// Create a knowledge graph with the mock driver
	kg := &KnowledgeGraph{driver: mockDriver}

	// Create a test state
	state := CursorState{
		Position: [3]float64{0.5, 0.6, 0.7},
		Mask:     [][][]float64{{{0.1, 0.2}, {0.3, 0.4}}},
	}

	// Diagnose the issue
	result, err := kg.DiagnoseIssue(context.Background(), state)
	assert.NoError(t, err)
	assert.Equal(t, "cursor_lag", result.Issue)
	assert.InDelta(t, 0.7, result.Confidence, 0.1)
	assert.Len(t, result.RelatedNodes, 3)
	assert.Contains(t, result.RelatedNodes, "hardware:mouse")
	assert.Contains(t, result.RelatedNodes, "software:driver")
	assert.Contains(t, result.RelatedNodes, "system:memory")

	// Verify mock expectations
	mockDriver.AssertExpectations(t)
	mockSession.AssertExpectations(t)
}

// TestKnowledgeGraphQueryFailure tests handling of Neo4j query failures
func TestKnowledgeGraphQueryFailure(t *testing.T) {
	// Create mock Neo4j driver and session
	mockDriver := new(MockNeo4jDriver)
	mockSession := new(MockNeo4jSession)

	// Set up expectations for a query failure
	mockDriver.On("NewSession", mock.Anything).Return(mockSession)
	mockSession.On("Run", mock.Anything, mock.Anything).Return(nil, assert.AnError)
	mockSession.On("Close").Return(nil)

	// Create a knowledge graph with the mock driver
	kg := &KnowledgeGraph{driver: mockDriver}

	// Create a test state
	state := CursorState{
		Position: [3]float64{0.5, 0.6, 0.7},
		Mask:     [][][]float64{{{0.1, 0.2}, {0.3, 0.4}}},
	}

	// Diagnose the issue
	_, err := kg.DiagnoseIssue(context.Background(), state)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to query knowledge graph")

	// Verify mock expectations
	mockDriver.AssertExpectations(t)
	mockSession.AssertExpectations(t)
}

// TestKnowledgeGraphNoResults tests handling of no results from Neo4j
func TestKnowledgeGraphNoResults(t *testing.T) {
	// Create mock Neo4j driver and session
	mockDriver := new(MockNeo4jDriver)
	mockSession := new(MockNeo4jSession)

	// Create empty result
	mockRecords := []neo4j.Record{}
	mockResult := NewMockNeo4jResult(mockRecords)

	// Set up expectations
	mockDriver.On("NewSession", mock.Anything).Return(mockSession)
	mockSession.On("Run", mock.Anything, mock.Anything).Return(mockResult, nil)
	mockSession.On("Close").Return(nil)
	mockResult.On("Err").Return(nil)
	mockResult.On("Consume").Return(nil, nil)

	// Create a knowledge graph with the mock driver
	kg := &KnowledgeGraph{driver: mockDriver}

	// Create a test state
	state := CursorState{
		Position: [3]float64{0.5, 0.6, 0.7},
		Mask:     [][][]float64{{{0.1, 0.2}, {0.3, 0.4}}},
	}

	// Diagnose the issue
	result, err := kg.DiagnoseIssue(context.Background(), state)
	assert.NoError(t, err)
	assert.Equal(t, "cursor_lag", result.Issue) // Default issue
	assert.Equal(t, 0.0, result.Confidence)     // Zero confidence
	assert.Empty(t, result.RelatedNodes)        // No related nodes

	// Verify mock expectations
	mockDriver.AssertExpectations(t)
	mockSession.AssertExpectations(t)
}

// TestKnowledgeGraphSessionError tests handling of Neo4j session errors
func TestKnowledgeGraphSessionError(t *testing.T) {
	// Create mock Neo4j driver
	mockDriver := new(MockNeo4jDriver)

	// Set up expectations for a session creation failure
	mockDriver.On("NewSession", mock.Anything).Return(nil)

	// Create a knowledge graph with the mock driver
	kg := &KnowledgeGraph{driver: mockDriver}

	// Create a test state
	state := CursorState{
		Position: [3]float64{0.5, 0.6, 0.7},
		Mask:     [][][]float64{{{0.1, 0.2}, {0.3, 0.4}}},
	}

	// Diagnose the issue
	_, err := kg.DiagnoseIssue(context.Background(), state)
	assert.Error(t, err)

	// Verify mock expectations
	mockDriver.AssertExpectations(t)
}
