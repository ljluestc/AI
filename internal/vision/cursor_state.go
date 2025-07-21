package vision

// EnhancedCursorState represents the state of the cursor with additional attributes
type EnhancedCursorState struct {
	Position     [3]float64    `json:"position"`
	Velocity     [3]float64    `json:"velocity,omitempty"`
	Acceleration [3]float64    `json:"acceleration,omitempty"`
	IsVisible    bool          `json:"isVisible"`
	Confidence   float64       `json:"confidence"`
	Mask         [][][]float64 `json:"mask,omitempty"`
}

// NewCursorState creates a new cursor state
func NewCursorState(position [3]float64, isVisible bool, confidence float64) *EnhancedCursorState {
	return &EnhancedCursorState{
		Position:   position,
		IsVisible:  isVisible,
		Confidence: confidence,
	}
}

// WithVelocity adds velocity to the cursor state
func (cs *EnhancedCursorState) WithVelocity(velocity [3]float64) *EnhancedCursorState {
	cs.Velocity = velocity
	return cs
}

// WithAcceleration adds acceleration to the cursor state
func (cs *EnhancedCursorState) WithAcceleration(acceleration [3]float64) *EnhancedCursorState {
	cs.Acceleration = acceleration
	return cs
}

// WithMask adds mask to the cursor state
func (cs *EnhancedCursorState) WithMask(mask [][][]float64) *EnhancedCursorState {
	cs.Mask = mask
	return cs
}
