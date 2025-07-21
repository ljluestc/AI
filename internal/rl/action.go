package rl

// Action represents an action that the agent can take
type Action struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Parameters  []float64 `json:"parameters,omitempty"`
}

// NewAction creates a new action with the given ID and description
func NewAction(id, description string) *Action {
	return &Action{
		ID:          id,
		Description: description,
		Parameters:  []float64{},
	}
}

// WithParameters adds parameters to the action
func (a *Action) WithParameters(params []float64) *Action {
	a.Parameters = params
	return a
}
