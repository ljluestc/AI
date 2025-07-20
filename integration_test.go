package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationEndToEnd tests the entire system flow from vision to execution
func TestIntegrationEndToEnd(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Create a test config
	config := Config{
		Port:          "8080",
		Neo4jURI:      "neo4j://localhost:7687",
		Neo4jUser:     "neo4j",
		Neo4jPassword: "password",
	}

	// Create a test agent
	agent, err := NewDebugAgent(config)
	assert.NoError(t, err)

	// Create a test router
	r := mux.NewRouter()
	r.HandleFunc("/vision", agent.handleVision).Methods("POST")
	r.HandleFunc("/diagnosis", agent.handleDiagnosis).Methods("POST")
	r.HandleFunc("/planning", agent.handlePlanning).Methods("POST")
	r.HandleFunc("/execution", agent.handleExecution).Methods("POST")

	// Step 1: Process vision input
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
	visionBody, _ := json.Marshal(visionInput)
	req, _ := http.NewRequest("POST", "/vision", bytes.NewBuffer(visionBody))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var cursorState CursorState
	err = json.Unmarshal(rr.Body.Bytes(), &cursorState)
	assert.NoError(t, err)
	assert.Len(t, cursorState.Position, 3)

	// Step 2: Diagnose the issue
	diagnosisBody, _ := json.Marshal(cursorState)
	req, _ = http.NewRequest("POST", "/diagnosis", bytes.NewBuffer(diagnosisBody))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	
	// If Neo4j is not available, we'll get an error, but we should continue the test
	var diagnosisResult DiagnosisResult
	if rr.Code == http.StatusOK {
		err = json.Unmarshal(rr.Body.Bytes(), &diagnosisResult)
		assert.NoError(t, err)
		assert.NotEmpty(t, diagnosisResult.Issue)
	}

	// Step 3: Plan an action
	planningInput := cursorState.Position[:]
	planningBody, _ := json.Marshal(planningInput)
	req, _ = http.NewRequest("POST", "/planning", bytes.NewBuffer(planningBody))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var action Action
	err = json.Unmarshal(rr.Body.Bytes(), &action)
	assert.NoError(t, err)
	assert.NotEmpty(t, action.ID)
	assert.NotEmpty(t, action.Description)

	// Step 4: Execute the action
	executionBody, _ := json.Marshal(action)
	req, _ = http.NewRequest("POST", "/execution", bytes.NewBuffer(executionBody))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var executionResult ExecutionResult
	err = json.Unmarshal(rr.Body.Bytes(), &executionResult)
	assert.NoError(t, err)
	assert.NotEmpty(t, executionResult.Message)
}

// TestAPIErrorHandling tests error handling in API endpoints
func TestAPIErrorHandling(t *testing.T) {
	// Create a test config
	config := Config{
		Port:          "8080",
		Neo4jURI:      "neo4j://localhost:7687",
		Neo4jUser:     "neo4j",
		Neo4jPassword: "password",
	}

	// Create a test agent
	agent, err := NewDebugAgent(config)
	assert.NoError(t, err)

	// Create a test router
	r := mux.NewRouter()
	r.HandleFunc("/vision", agent.handleVision).Methods("POST")
	r.HandleFunc("/diagnosis", agent.handleDiagnosis).Methods("POST")
	r.HandleFunc("/planning", agent.handlePlanning).Methods("POST")
	r.HandleFunc("/execution", agent.handleExecution).Methods("POST")

	// Test vision endpoint with invalid input
	req, _ := http.NewRequest("POST", "/vision", bytes.NewBufferString("invalid json"))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid input")

	// Test diagnosis endpoint with invalid input
	req, _ = http.NewRequest("POST", "/diagnosis", bytes.NewBufferString("invalid json"))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid state")

	// Test planning endpoint with invalid input
	req, _ = http.NewRequest("POST", "/planning", bytes.NewBufferString("invalid json"))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid state")

	// Test execution endpoint with invalid input
	req, _ = http.NewRequest("POST", "/execution", bytes.NewBufferString("invalid json"))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid action")
}

// TestWebSocketHandling tests the WebSocket connection
func TestWebSocketHandling(t *testing.T) {
	// Skip in short mode or if WebSocket testing is not available
	t.Skip("Skipping WebSocket test - requires special setup")
	
	// Note: WebSocket testing requires more complex setup with a real server
	// This is a placeholder for a WebSocket test
}
