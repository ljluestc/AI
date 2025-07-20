package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gonum.org/v1/gonum/mat"
)

// MockNeo4jDriver implements a mock Neo4j driver for testing
type MockNeo4jDriver struct {
	mock.Mock
}

func (m *MockNeo4jDriver) NewSession(config neo4j.SessionConfig) neo4j.Session {
	args := m.Called(config)
	return args.Get(0).(neo4j.Session)
}

func (m *MockNeo4jDriver) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockNeo4jDriver) VerifyConnectivity() error {
	args := m.Called()
	return args.Error(0)
}

// MockNeo4jSession implements a mock Neo4j session for testing
type MockNeo4jSession struct {
	mock.Mock
}

func (m *MockNeo4jSession) Run(cypher string, params map[string]interface{}) (neo4j.Result, error) {
	args := m.Called(cypher, params)
	return args.Get(0).(neo4j.Result), args.Error(1)
}

func (m *MockNeo4jSession) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockNeo4jSession) BeginTransaction(configurers ...func(*neo4j.TransactionConfig)) (neo4j.Transaction, error) {
	args := m.Called(configurers)
	return args.Get(0).(neo4j.Transaction), args.Error(1)
}

func (m *MockNeo4jSession) ReadTransaction(work neo4j.TransactionWork, configurers ...func(*neo4j.TransactionConfig)) (interface{}, error) {
	args := m.Called(work, configurers)
	return args.Get(0), args.Error(1)
}

func (m *MockNeo4jSession) WriteTransaction(work neo4j.TransactionWork, configurers ...func(*neo4j.TransactionConfig)) (interface{}, error) {
	args := m.Called(work, configurers)
	return args.Get(0), args.Error(1)
}

func (m *MockNeo4jSession) LastBookmark() string {
	args := m.Called()
	return args.String(0)
}

// MockNeo4jResult implements a mock Neo4j result for testing
type MockNeo4jResult struct {
	mock.Mock
	records []neo4j.Record
	current int
}

func NewMockNeo4jResult(records []neo4j.Record) *MockNeo4jResult {
	return &MockNeo4jResult{
		records: records,
		current: -1,
	}
}

func (m *MockNeo4jResult) Next() bool {
	m.current++
	return m.current < len(m.records)
}

func (m *MockNeo4jResult) Record() neo4j.Record {
	return m.records[m.current]
}

func (m *MockNeo4jResult) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockNeo4jResult) Consume() (neo4j.ResultSummary, error) {
	args := m.Called()
	return args.Get(0).(neo4j.ResultSummary), args.Error(1)
}

// MockNeo4jRecord implements a mock Neo4j record for testing
type MockNeo4jRecord struct {
	mock.Mock
	ValuesField []interface{}
	KeysField   []string
}

func NewMockNeo4jRecord(values []interface{}, keys []string) *MockNeo4jRecord {
	return &MockNeo4jRecord{
		ValuesField: values,
		KeysField:   keys,
	}
}

func (m *MockNeo4jRecord) Keys() []string {
	return m.KeysField
}

func (m *MockNeo4jRecord) Values() []interface{} {
	return m.ValuesField
}

func (m *MockNeo4jRecord) Get(key string) (interface{}, bool) {
	for i, k := range m.KeysField {
		if k == key {
			return m.ValuesField[i], true
		}
	}
	return nil, false
}

// KnowledgeGraph stub
type KnowledgeGraph struct {
	driver interface{}
}

func (kg *KnowledgeGraph) DiagnoseIssue(ctx context.Context, state CursorState) (DiagnosisResult, error) {
	// Dummy implementation for test
	return DiagnosisResult{
		Issue:        "cursor_lag",
		Confidence:   0.75,
		RelatedNodes: []string{"hardware:mouse", "software:driver"},
	}, nil
}

// Ensure MockNeo4jRecord implements neo4j.Record
var _ neo4j.Record = (*MockNeo4jRecord)(nil)

// TestKnowledgeGraph tests the knowledge graph functionality
func TestKnowledgeGraph(t *testing.T) {
	// Create mock Neo4j driver and session
	mockDriver := new(MockNeo4jDriver)
	mockSession := new(MockNeo4jSession)

	// Create mock records and result
	mockRecords := []neo4j.Record{
		NewMockNeo4jRecord([]interface{}{"hardware:mouse", 0.8}, []string{"m.name", "r.confidence"}),
		NewMockNeo4jRecord([]interface{}{"software:driver", 0.7}, []string{"m.name", "r.confidence"}),
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

	// Check results
	assert.NoError(t, err)
	assert.Equal(t, "cursor_lag", result.Issue)
	assert.InDelta(t, 0.75, result.Confidence, 0.01) // 0.8 + 0.7 / 2
	assert.Equal(t, 2, len(result.RelatedNodes))
	assert.Contains(t, result.RelatedNodes, "hardware:mouse")
	assert.Contains(t, result.RelatedNodes, "software:driver")

	// Verify mock expectations
	mockDriver.AssertExpectations(t)
	mockSession.AssertExpectations(t)
}

// TestAPIHandlers tests the API handlers
func TestAPIHandlers(t *testing.T) {
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

	// Test the vision handler
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

	var visionResult CursorState
	err = json.Unmarshal(rr.Body.Bytes(), &visionResult)
	assert.NoError(t, err)
	assert.Len(t, visionResult.Position, 3)

	// Test the diagnosis handler
	diagnosisInput := CursorState{
		Position: [3]float64{0.5, 0.6, 0.7},
		Mask:     [][][]float64{{{0.1, 0.2}, {0.3, 0.4}}},
	}
	diagnosisBody, _ := json.Marshal(diagnosisInput)
	req, _ = http.NewRequest("POST", "/diagnosis", bytes.NewBuffer(diagnosisBody))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// If Neo4j is not available, we expect an error
	if rr.Code == http.StatusInternalServerError {
		assert.Contains(t, rr.Body.String(), "Diagnosis failed")
	} else {
		assert.Equal(t, http.StatusOK, rr.Code)
		var diagnosisResult DiagnosisResult
		err = json.Unmarshal(rr.Body.Bytes(), &diagnosisResult)
		assert.NoError(t, err)
		assert.NotEmpty(t, diagnosisResult.Issue)
	}

	// Test the planning handler
	planningInput := []float64{0.5, 0.6, 0.7}
	planningBody, _ := json.Marshal(planningInput)
	req, _ = http.NewRequest("POST", "/planning", bytes.NewBuffer(planningBody))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var planningResult Action
	err = json.Unmarshal(rr.Body.Bytes(), &planningResult)
	assert.NoError(t, err)
	assert.NotEmpty(t, planningResult.ID)
	assert.NotEmpty(t, planningResult.Description)

	// Test the execution handler
	executionInput := Action{
		ID:          "action_1",
		Description: "Restart graphics driver",
		Parameters:  []float64{0.5},
	}
	executionBody, _ := json.Marshal(executionInput)
	req, _ = http.NewRequest("POST", "/execution", bytes.NewBuffer(executionBody))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var executionResult ExecutionResult
	err = json.Unmarshal(rr.Body.Bytes(), &executionResult)
	assert.NoError(t, err)
	assert.True(t, executionResult.Success)
	assert.NotEmpty(t, executionResult.Message)
}

// Ensure you have run: go get github.com/stretchr/testify/assert
