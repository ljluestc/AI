package knowledge

// MockNeo4jRecord is a mock implementation of neo4j.Record for testing
type MockNeo4jRecord struct {
	values []interface{}
	keys   []string
}

// NewMockNeo4jRecord creates a new MockNeo4jRecord
func NewMockNeo4jRecord(values []interface{}, keys []string) *MockNeo4jRecord {
	return &MockNeo4jRecord{
		values: values,
		keys:   keys,
	}
}

// Keys returns the keys of the record
func (r *MockNeo4jRecord) Keys() []string {
	return r.keys
}

// Values returns the values of the record
func (r *MockNeo4jRecord) Values() []interface{} {
	return r.values
}

// Get returns the value for the given key
func (r *MockNeo4jRecord) Get(key string) (interface{}, bool) {
	for i, k := range r.keys {
		if k == key {
			return r.values[i], true
		}
	}
	return nil, false
}
